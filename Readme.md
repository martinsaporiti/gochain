# Blockchain sample coded with Go :rocket:
This repository contains the code to run a blockchain solution and a wallet client.
I coded this to learn how a blockchain works and get better understanding about public and private keys, and consensus algorithms.
Some parts are solved trying to keep things simples, for instance, nothing is saved in a database. The blockchain and the wallet keeps data only in memory.

With this code you will be able to:
* Create user's wallets.
* Run n nodes, each one with its own version of the blockchain.


### How run the nodes
For instance, for running 3 nodes you should execute:
```bash
# This runs node 1
go run cmd/blockchain/main.go -port 5000
```

```bash
# This runs node 2
go run cmd/blockchain/main.go -port 5001
```

```bash
# This runs node 3
go run cmd/blockchain/main.go -port 5002
```
The ports numbers are **important** because every node looks up neighbords in a range of ips and ports.
You can see that taking a look in `./internal/gateway/http_gateway.go`:

```go
const (
	BLOCKCHAIN_PORT_RANGE_START      = 5000
	BLOCKCHAIN_PORT_RANGE_END        = 5003
	NEIGHBOR_IP_RANGE_START          = 1
	NEIGHBOR_IP_RANGE_END            = 3
	BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC = 10
)
```
For running the wallet client, you must run:

```bash
go run cmd/wallet/main.go -port 8080 -gateway "http://localhost:5000"
```

The `gateway` parameter must be one of the nodes created before. Then you can see the wallet visiting `http://localhost:8080/`

**Important**:
Take in mind every time you reload the page a new wallet is created (a new blockchain address, private key, and public key). To send money from one wallet to another, you must open two tabs in your browser, creating two wallets.
## How to see the blockchain
You can see the blockchain calling:
```bash
# This is for node 5000
http://localhost:5000/
```
and you should get something like this after sending USD 100 from one wallet to another:
```json
{
	"chain": [
		{
			"number": 1,
			"nonce": 0,
			"previous_hash": "c1e7b244fb428cff0cc72011321454a080ab41f189...",
			"timestamp": 1654695626111823000,
			"transactions": null
		},
		{
			"number": 2,
			"nonce": 2916576,
			"previous_hash": "e0a039936b65f36ff6e3dfdc9c03...",
			"timestamp": 1654695720103761000,
			"transactions": [
				{
					"sender_blockchain_address": "18fwCkKmcPJonyScY7qqThgbg1WPVd2aA1",
					"recipient_blockchain_address": "1FsRTaZ2LoPafdjMr9qwnkyPEkn5jDB6dk",
					"value": 100,
					"timestamp": 1654695654
				},
				{
					"sender_blockchain_address": "THE BLOCKCHAIN 5000",
					"recipient_blockchain_address": "136KiUxSRZg2padBDmdkH4oqDb51F3TiKi",
					"value": 1,
					"timestamp": 1654695659
				}
			]
		}
	]
}
```
if you have more nodes running you can call `http://localhost:500X/`
