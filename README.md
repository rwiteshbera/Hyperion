![Hyperion Header](https://user-images.githubusercontent.com/73098407/210416970-710f9153-18fe-48fd-9e51-39db3b9801ed.jpg)

#### Blockchain project aims to implement the core functionality of a blockchain, including the ability to add new blocks to the chain, verify the integrity of the chain, and secure the data stored in the blocks using cryptographic hashes.

- Create New Blockchain Wallet
```bash
go run main.go --create
```
**Warning:** Never disclose the private key. Anyone with your private key can steal your assets. 


- List All Blockchain Wallets
```bash
go run main.go --list
```

- Get Private Key and Public Key for Transaction
```bash
go run main.go --get <WALLET_ADDRESS>
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

- Check Wallet Balance
```text
Endpoint : http://localhost:<PORT>/wallet/<WALLET_ADDRESS>
Request Type : GET
```
