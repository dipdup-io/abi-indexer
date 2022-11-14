# ABI Indexer
Ethereum ABI indexer receives ABI from third-party sources and stores it to database. Now indexer receives ABI from Sourcify or file system (via volume in Docker case).

## Install

```bash
make build

# or

docker-compose up -d --build
```

## Configuration

Example of config file can be found [here](/build/dipdup.yml)

Environment variables:

```bash
LOG_LEVEL=debug                           # logging level: trace, debug, info, warning, error 
VM_TYPE=evm                               # using virtual machine type: now supported only EVM
METADATA_SOURCE_TYPE=sourcify             # source of ABI: sourcify or fs
METADATA_THREADS_COUNT=10                 # receiving workers count
SOURCIFY_BASE_URL=https://sourcify.dev    # Sourcify base URL
SOURCIFY_CHAIN_ID=1                       # Sourcify chain ID. Can be found here: https://sourcify.dev/server/chains
SOURCIFY_TIMEOUT=10                       # timeout request to Sourcify
FS_DIR=/etc/metadata                      # directory which used for File System source of ABI
GRPC_BIND=127.0.0.1:7778                  # which hostname:port will be used for gRPC
POSTGRES_HOST=127.0.0.1                   # Postgres hostname
POSTGRES_PORT=5432                        # Postgres port
```

## API

You can communicate with module by [gRPC](/pkg/modules/grpc).