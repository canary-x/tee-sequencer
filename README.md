# TEE sequencer

Responsible for shuffling sequences of Ethereum transactions in a nonce-honoring way.
Set up to run in an AWS Nitro Enclave.

## TODO
- add a linter
- find a way to expose logs to the host
- add support for SIGINT and SIGTERM signals
- replace with efficient protobuf implementation

## Building the Nitro enclave
On an EC2 instance with the nitro-cli already installed:
```shell
make docker/build
nitro-cli build-enclave --docker-uri com.github.canary-x.tee-sequencer:latest --output-file sequencer.eif
```

The output should be as follows:
```json
{
  "Measurements": {
    "HashAlgorithm": "Sha384 { ... }",
    "PCR0": "a2d89ca2c8f451fa469a67c4282a9be33375802deb6264e82ce93a8551fe451bcdf90c69661570844f40696964a24e0c",
    "PCR1": "4b4d5b3661b3efc12920900c80e126e4ce783c522de6c02a2a5bf7af3a2b9327b86776f188e4be1c1c404a129dbda493",
    "PCR2": "b9201487cba85799674dc2df05ba9ec34bf3d7779dd1a2c8c7e491b8395d58d6265a0ef57735d4a860f2b1cdf261805a"
  }
}
```

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
    "Measurements": {
      "HashAlgorithm": "Sha384 { ... }",
      "PCR0": "a2d89ca2c8f451fa469a67c4282a9be33375802deb6264e82ce93a8551fe451bcdf90c69661570844f40696964a24e0c",
      "PCR1": "4b4d5b3661b3efc12920900c80e126e4ce783c522de6c02a2a5bf7af3a2b9327b86776f188e4be1c1c404a129dbda493",
      "PCR2": "b9201487cba85799674dc2df05ba9ec34bf3d7779dd1a2c8c7e491b8395d58d6265a0ef57735d4a860f2b1cdf261805a"
    }
  }
```
Notice the `CID` being `16`, which we'll use as a vsock ID for connections.
You can check the logs while in debug mode: ```nitro-cli console --enclave-name sequencer```

### Submitting test transactions
You can run an HTTP proxy to forward requests to the enclave:
```shell
export VSOCK=$(nitro-cli describe-enclaves | jq '.[] | select(.EnclaveName == "sequencer") | .EnclaveCID')
docker run -d -p 8080:8080 --name socat alpine/socat tcp-listen:8080,fork,reuseaddr vsock-connect:$VSOCK:8080
```

Then you can submit test transactions to the enclave:
```shell
curl --location --request GET 'http://localhost:8080' \
--header 'Content-Type: application/json' \
--data '{
    "transactions": [
        {
            "tx_hash": "hash-1",
            "account": "account-1",
            "nonce": 1
        },
        {
            "tx_hash": "hash-1",
            "account": "account-1",
            "nonce": 2
        }
    ]
}'
```

Cleanup:
```shell
nitro-cli terminate-enclave --enclave-id $(nitro-cli describe-enclaves | jq -r '.[] | select(.EnclaveName == "sequencer") | .EnclaveID')
docker rm -f socat
```