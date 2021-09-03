# eth-log-aggregator

## Usage

The usage of this tool is exactly the same as [`abigen`](https://geth.ethereum.org/docs/dapp/native-bindings). It is expected
that bindings for a contract exist before running this tool.

Once the tool has been built `go build`, you can simply generate the log aggregate bindings as follows:

```console
abigen --abi contracts/ERC20.json --pkg contracts --type Token --out pkg/contracts/token.go
eth-log-aggregator --abi contracts/ERC20.json --pkg contracts --type Token --out pkg/contracts/token_watcher.go
```
