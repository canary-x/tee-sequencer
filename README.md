# TEE sequencer

Responsible for shuffling sequences of Ethereum transactions in a nonce-honoring way.
Set up to run in an AWS Nitro Enclave.

## Work in progress

- add a linter
- find a way to expose logs to the host
- add support for SIGINT and SIGTERM signals
- add signature to response
- implement actual shuffling

## Development

Nitro enclaves only allow vsock networking, with Linux being the only OS supporting it.
When developing on other OSes, you can either use docker or simply run the sequencer as is, which will detect the lack
of vsock support and fall back to a regular TCP socket.

### Dependencies

- go 1.21 or later
- make
- buf (install via `make deps`)

## Building the Nitro enclave

### Prerequisites

1. An AWS m5.xlarge instance with Amazon Linux 2023 (not strictly necessary, but useful to reproduce these steps with
   confidence)
2. Follow the AWS Nitro getting started [guide](https://docs.aws.amazon.com/enclaves/latest/user/getting-started.html)
3. Install deps: ```yum install -y socat make jq```

### Build

```shell
make docker/build
nitro-cli build-enclave --docker-uri com.github.canary-x.tee-sequencer:latest --output-file sequencer.eif
```

The output should be as follows:

```json
{
  "Measurements": {
    "HashAlgorithm": "Sha384 { ... }",
    "PCR0": "1d91531c44241c530d4e6cdab913d5ca348d0922bd13a3b26ce75edf0c249707b38b1a53ac39461e79ae82c483e695ee",
    "PCR1": "...(this depends on your host instance)...",
    "PCR2": "a82e249a0453597c949a0ed2a2b223f1febd61c320a265036eb61bf5a3397d2603e19d7cb930ccf8d5bb2ee720fd9c13"
  }
}
```

### Run

Running in debug mode, assuming 2 CPU cores were allocated and at least 512MiB RAM:

```shell
nitro-cli run-enclave --cpu-count 2 --memory 512 --eif-path sequencer.eif --debug-mode
nitro-cli describe-enclaves
```

The output should be something like:

```json
  {
  "EnclaveName": "sequencer",
  "EnclaveID": "i-0bdbf2d4b4b7b2e35-enc191fb82b45f290e",
  "ProcessID": 4519,
  "EnclaveCID": 16,
  "NumberOfCPUs": 2,
  "CPUIDs": [
    1,
    3
  ],
  "MemoryMiB": 512,
  "State": "RUNNING",
  "Flags": "DEBUG_MODE",
  "Measurements": "..."
}
```

Notice the `CID` being `16`, which we'll use as a vsock ID for connections.
You can check the logs while in debug mode: ```nitro-cli console --enclave-name sequencer```

### Submitting test transactions

You can run an HTTP proxy to forward requests to the enclave:

```shell
export VSOCK=$(nitro-cli describe-enclaves | jq '.[] | select(.EnclaveName == "sequencer") | .EnclaveCID')
socat TCP-LISTEN:8080,reuseaddr,fork VSOCK-CONNECT:$VSOCK:8080 &
```

Then you can submit test transactions to the enclave:

```shell
curl --location 'http://localhost:8080/blockchain.v1.SequencerService/Shuffle' \
--header 'Content-Type: application/json' \
--data '{
    "transactions": [
        {
            "tx_hash": "aGFzaC0x",
            "account": "YWNjb3VudC0x",
            "nonce": "bm9uY2Ux"
        },
        {
            "tx_hash": "aGFzaC0y",
            "account": "YWNjb3VudC0y",
            "nonce": "bm9uY2Uy"
        }
    ]
}'
```

Note: protobuf byte fields are represented as base64 strings when using application/json encoding.

Cleanup:

```shell
nitro-cli terminate-enclave --enclave-id $(nitro-cli describe-enclaves | jq -r '.[] | select(.EnclaveName == "sequencer") | .EnclaveID')
ps aux | grep socat | grep -v grep | awk '{print $2}' | xargs -r sudo kill -9
```

### Collecting logs

A special zap logger will stream logs to both the console and a vsock connection.
Console logs are only visible in DEBUG mode, which is why there's a necessity for an additional stream.
In order to collect logs, run the special [nitro-logger](https://github.com/canary-x/nitro-logger) as a service on your
parent instance, on port 9000.

## Documentation

Curious about how this tech works? Check out our in-depth [documentation](DOC.md).