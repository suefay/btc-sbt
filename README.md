## BTC-SBT

The [BTC-SBT](https://btc-sbt.gitbook.io/btc-sbt-protocol) protocol implementation along with the Command Line Interface.

## Getting Started

### Configuration

The config file is required to use the CLI.

The default config file is `config.yaml` in the working directory.

The items to be configured are as follows:

- node
  - rpc_url: node rpc url

  - rpc_user: node rpc username

  - rpc_pass: node rpc passphrase

  - net_version: net version(0:mainnet, 1:testnet, 2:signet)

- unisat:
  - api: unisat api

- indexer:
  - interval: indexer interval

- db
  - path: db path

- key_store
  - path: key file in which the key WIF is stored

- fee_rate: fee rate for initiating txs

- server
  - listener_address: listener address for the server

- general
  - retries: retry count

  - interval: retry interval

- log
  - level: logging level(1:fatal, 2:error, 3:warning, 4:info, 5:debug, 6:trace)

### Start BTC-SBT node

```bash
btc-sbt node [config-file]
```

### Issue BTC SBT

```bash
btc-sbt issue [args] [config-file]
```

### Mint BTC SBT

```bash
btc-sbt mint [args] [config-file]
```
