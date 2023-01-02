Blockchain project aims to implement the core functionality of a blockchain, including the ability to add new blocks to the chain, verify the integrity of the chain, and secure the data stored in the blocks using cryptographic hashes.

- Create New Blockchain Wallet
```bash
go run main.go --create
```
Keep your private key secure and do not share it. It is stored in encrypted form on your local device and can only be accessed with the correct passphrase.

- List All Blockchain Wallets
```bash
go run main.go --list
```

- Run Blockchain Server
```bash
go run main.go --port <PORT_NUMBER>
```
- Visualize Blocks in Hyperion
```text
Endpoint : http://localhost:<PORT>
Request Type : GET
```
- Make new transaction
```text
Endpoint : http://localhost:<PORT>/new
Request Type : POST

Request Body in JSON Format:
{
"privateKey" : "",
"publicKey" : "",
"sender" : "",
"recipient" : "",
"value" : 1
}
```
- Verify Transaction
```text
Endpoint : http://localhost:<PORT>/explore/<TRANSACTION_HASH>
Request Type : GET
```