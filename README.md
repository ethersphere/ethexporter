# ETHexporter

This project is a fork from [hunterlong/ethexporter](https://github.com/hunterlong/ethexporter)

A lightweight Prometheus exporter and [Grafana Dashboard](https://grafana.com/dashboards/6970) that will output Ethereum wallet balances from a list of addresses you specify. ETHexporter attaches to a geth server to fetch ethereum wallet balances for your Grafana dashboards. You can also use [TOKENexporter](https://github.com/hunterlong/tokenexporter) for ERC20 token balances.

## watch addresses
The `addresses.txt` file holds all the addresses to fetch balances for. Use the format `name:address` on each new line.
```
etherdelta:0x8d12A197cB00D4747a1fe03395095ce2A5CC6819
bittrex:0xFBb1b73C4f0BDa4f67dcA266ce6Ef42f520fBB98
```

## grafana dashboard
ETHexporter includes a [Grafana Dashboard](https://grafana.com/dashboards/6970) so you visualize Ethereum wallet balances with ease. Import Grafana Dashboard ID: `6970` once you've started your
Docker container.

<p align="center"><img width="70%" src="https://img.cjx.io/ethexportergrafana.png"></p>

## build docker image
Clone this repo and then follow the simple steps below!

### build docker image
`docker build -t ethexporter:latest .`

### run ethexporter
You'll need access to an ethereum Geth server to fetch balances. You can use [Infura.io](https://infura.io/setup) or [Slock.it](https://rpc.slock.it/).
```
docker run -d -p 9890:9890 \
 -e GETH="https://rpc.slock.it/mainnet" \
 ethersphere/ethexporter:latest
```

## pull from dockerhub
Create a `addresses.txt` file with the correct format mentioned above.
```
docker run -d -v $(pwd)/addresses.txt:/app/addresses.txt \
 -p 9890:9890 \
 -e GETH=https://rpc.slock.it/goerli \
 ethersphere/ethexporter:latest
```
The Docker image should be running with the default addresses.

## /metrics response
```
eth_balance{name="etherdelta",address="0x8d12A197cB00D4747a1fe03395095ce2A5CC6819"} 24919.37437
eth_balance{name="bittrex",address="0xFBb1b73C4f0BDa4f67dcA266ce6Ef42f520fBB98"} 687509.5097
eth_balance{name="poloniex",address="0x32Be343B94f860124dC4fEe278FDCBD38C102D88"} 72284.47401
eth_balance{name="kraken",address="0x267be1c1d684f78cb4f6a176c4911b741e4ffdc0"} 159592.0022
eth_balance_total 944305.360280000022612512
eth_load_seconds 1.15
eth_loaded_addresses 4
eth_total_addresses 4
```
