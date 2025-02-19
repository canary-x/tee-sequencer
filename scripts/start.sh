#!/usr/bin/env bash
set -e

# Start or restart the sequencer enclave and the socat proxy

cpu_count=2
memory=512
debug_mode=false
eif_path="sequencer.eif"
port=8080

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --cpu-count)
            cpu_count="$2"
            shift 2
            ;;
        --memory)
            memory="$2"
            shift 2
            ;;
        --debug-mode)
            debug_mode=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

if ! command -v socat &> /dev/null; then
    echo "Error: socat is not installed. Please install it and try again." >&2
    exit 1
fi

if ! command -v nitro-cli &> /dev/null; then
    echo "Error: nitro-cli is not installed. Please install it and try again." >&2
    exit 1
fi

enclave_info=$(nitro-cli describe-enclaves)
enclave_id=$(echo "$enclave_info" | jq -r '.[0].EnclaveID')

if [ -n "$enclave_id" ]; then
    echo "Sequencer is already running with ID $enclave_id. Terminating..."
    nitro-cli terminate-enclave --enclave-id "$enclave_id"
fi

echo "Starting sequencer..."
cmd="nitro-cli run-enclave --cpu-count $cpu_count --memory $memory --eif-path $eif_path"
if [ "$debug_mode" = true ]; then
    cmd+=" --debug-mode"
fi

echo "Executing command: $cmd"
eval "$cmd"

echo "Waiting for enclave to start running..."
sleep 1

check_enclave_running() {
    local timeout=30
    local interval=3
    local start_time=$(date +%s)

    while true; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))

        if [ $elapsed -ge $timeout ]; then
            echo "Timeout reached. Enclave did not start within $timeout seconds." >&2
            return 1
        fi

        local enclave_info=$(nitro-cli describe-enclaves)

        if [ -z "$enclave_info" ] || [ "$enclave_info" == "[]" ]; then
            echo "Enclave not found. It may have terminated with an error." >&2
            return 1
        fi

        local state=$(echo "$enclave_info" | jq -r '.[0].State')

        if [ "$state" == "RUNNING" ]; then
            echo "Enclave is now running."
            echo "$enclave_info" | jq '.[0]'
            return 0
        else
            echo "Waiting for enclave to start. Current state: $state" >&2
            sleep $interval
        fi
    done
}

if ! check_enclave_running; then
    echo "Failed to start the enclave. Exiting."
    exit 1
fi

port_check=$(sudo lsof -i :$port -P -n -sTCP:LISTEN)
if [ -n "$port_check" ]; then
    process_name=$(echo "$port_check" | tail -n 1 | awk '{print $1}')
    pid=$(echo "$port_check" | tail -n 1 | awk '{print $2}')
    if [ "$process_name" = "socat" ]; then
        echo "Killing existing socat process on port $port"
        sudo kill $pid
    else
        echo "Error: Port $port is occupied by $process_name (PID: $pid)"
        exit 1
    fi
fi

vsock=$(nitro-cli describe-enclaves | jq '.[] | select(.EnclaveName == "sequencer") | .EnclaveCID')
echo "Starting socat on CID:port $vsock:$port"
socat TCP-LISTEN:$port,reuseaddr,fork vsock-connect:$vsock:$port &
