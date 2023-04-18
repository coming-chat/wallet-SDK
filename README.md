# wallet-SDK

|                       | Bitcoin | Ethereum | Polka | Cosmos | Doge | Solana | Aptos | Sui | Starcoin |
| --------------------- | ------- | -------- | ----- | ------ | ------ | ------ | ------ | ------ | ------ |
| import mnemonic       | ✅     | ✅      | ✅   |✅   |✅   |✅   |✅   |✅   |✅   |
| import keystore       | ❌     | ❌      | ✅  |❌   |❌   |❌   |❌   |❌   |❌   |
| pri/pub key & address | ✅     | ✅      | ✅   |✅   |✅   |✅   |✅   |✅   |✅   |
| multi network         | ✅     | ✅      | ✅   |✅   |❌   |✅   |✅   |✅   |✅   |
| publicKey to address  | ✅     | ✅      | ✅   |✅   |✅   |✅   |✅   |✅   |✅   |
| address to publicKey  | ❌     | ❌      | ✅   |❌   |❌   |✅   |❌   |❌   |❌   |
| sign data             | ☑️    | ✅    | ✅   |☑️   |☑️   |☑️   |✅   |✅   |✅   |
|                       |         |          |       |        |      |        |        |        |        |
| query balance | ✅ | ✅ | ✅ |✅   |✅   |✅   |✅   |✅   | ✅ |
| fetch transaction detail | ✅ | ✅ | ✅ |✅   |✅   |✅   |✅   |✅   | ✅ |
| gas fee | ❌ | ✅ | ✅ |☑️   |✅   |✅   |✅   |❌   | ✅ |
| send raw transaction | ✅ | ✅ | ✅ ☑️ |✅   |✅   |✅   |✅   |✅   | ✅ |
| multi token | ❌ | ✅ erc20 | ✅ XBTC |✅   |❌   |❌   |✅   |❌   | ❌ |

*If there are two icons, the first icon indicates development status, and the second icon indicates test status.*
✅: Completed      ☑️: TODO    ❌: Unsupported

## Usage

### About Wallet 

SDK provide wallet import, account public and private key and address acquisition.

#### Import Wallet

```golang
// import mnemonic
wallet, err = NewWalletFromMnemonic(mnemonic)

// import keystore
// It only supports Polka keystore.
wallet, err = NewWalletFromKeyStore(keyStoreJson, password)
```

#### Create Account

We currently support accounts in Bitcoin, Ethereum and Polkadot ecosystems.

```golang
// Polka
polkaAccount, err = wallet.GetOrCreatePolkaAccount(network)

// Bitcoin
bitcoinAccount, err = wallet.GetOrCreateBitcoinAccount(chainnet)

// Ethereum
ethereumAccount, err = wallet.GetOrCreateEthereumAccount()

// Cosmos
cosmosAccount, err = wallet.GetOrCreateCosmosAccount()

// Terra is a type of cosmos
terraAccount, err = wallet.GetOrCreateCosmosTypeAccount(330, "terra")
```

#### Get PrivateKey, PublicKey, Address

```golang
privateData, err = account.PrivateKeyData()

privateKey, err = account.PrivateKey()

publicKey = account.PublicKey()

address = account.Address()
```


### About Chain

We can use chain tools to do chain related work.
* query balance
* query estimate fees
* send transaction
* fetch transaction detail
* support multi token: 
  * eth contract erc20 token
  * xbtc

#### Create Chain

```go
polkaChain, err = polka.NewChainWithRpc(rpcUrl, scanUrl, network)

bitcoinChain, err = btc.NewChainWithChainnet(chainnet)

ethereumChain, err = eth.NewChainWithRpc(rpcUrl)

cosmosChain, err = cosmos.NewChainWithRpc(rpcUrl, restUrl)
terraChain, err = cosmos.NewChainWithRpc(terraRpcUrl, terraRestUrl)
```

#### Methods

```golang 

// query balance
balance, err = chain.BalanceOfAddress(address)
balance, err = chain.BalanceOfPublicKey(publicKey)
balance, err = chain.BalanceOfAccount(account)

// send transaction
txHash, err = chain.SendRawTransaction(signedTx)

// fetch transaction detail
detail, err = chain.FetchTransactionDetail(hashString)
status = chain.FetchTransactionStatus(hashString)
```

#### Chain's Token

```golang
// MainToken
token = chain.MainToken()

// btc have not tokens

// polka (only support XBTC of ChainX currently)
xbtcToken = polkaChain.XBTCToken()

// ethereum erc20 token
erc20Token = ethereumChain.Erc20Token(contractAddress)

// cosmos token
atomToken = cosmosChain.DenomToken("cosmos", "uatom")
ustToken = terraChain.DenomToken("terra", "uusd")

// token balance (similar to chain's balance)
balance, err = anyToken.BalanceOfAddress(address)
balance, err = anyToken.BalanceOfPublicKey(publicKey)
balance, err = anyToken.BalanceOfAccount(account)
```

#### Estimate fee

```golang
// sbtc's estimate fee is compute by utxo

// polka 
transaction = // ...
fee, err = polkaChain.EstimateFeeForTransaction(transaction)

// ethereum
gasPrice, err = ethChain.SuggestGasPrice()
gasLimit, err = anyEthToken.EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount)
```

--------------------------------------------------------------------------------
--------------------------------------------------------------------------------



## Cosmos & Terra Demo

### Creator

```go
// Import mnemonic
wallet, err = NewWalletFromMnemonic(mnemonic)

// Cosmos account
cosmosAccount, err = wallet.GetOrCreateCosmosAccount()
// Terra is a type of cosmos
terraAccount, err = wallet.GetOrCreateCosmosTypeAccount(330, "terra")

// Create Chain
chain, err = cosmos.NewChainWithRpc(rpcUrl, restUrl)

// Create a coin token
atomToken = cosmosChain.DenomToken("cosmos", "uatom")
lunaToken = terraChain.DenomToken("terra", "uluna")
ustToken = terraChain.DenomToken("terra", "uusd")
```

### Query balance

```go
// query balance with token
atomBalance, err = atomToken.BalanceOfAddress("cosmos1lkw6n8efpj7mk29yvajpn9zue099l359cgzf0t")
lunaBalance, err = lunaToken.BalanceOfAddress("terra1ncjg4a59x2pgvqy9qjyqprlj8lrwshm0wleht5")
ustBalance , err =  ustToken.BalanceOfAddress("terra1dr7ackrxsqwmac2arx26gre6rj6q3sv29fnn7k")

// you can query balance with chain directly
atomBalance, err = cosmosChain.BalanceOfAddressAndDenom("cosmos1lkw6n8efpj7mk29yvajpn9zue099l359cgzf0t", "uatom") // uatom
lunaBalance, err =  terraChain.BalanceOfAddressAndDenom("terra1ncjg4a59x2pgvqy9qjyqprlj8lrwshm0wleht5", "uluna")  // uluna
ustBalance , err =  terraChain.BalanceOfAddressAndDenom("terra1dr7ackrxsqwmac2arx26gre6rj6q3sv29fnn7k", "uusd")   // uusd
```

### Fetch Transaction Detail

```go
atomDetail, err = cosmosChain.FetchTransactionDetail("F068275DE4A4CC904D3E6A412A50DFACC235C62770BCD001E54E00BC4C17B1F0")

lunaDetail, err =  terraChain.FetchTransactionDetail("19771A22934641DBD3D347DCCAE939DAC37F39ABD88005AA735B8AAEA78599BA")

// exactly the same as luna detail.
ustDetail , err =  terraChain.FetchTransactionDetail("25ACF16526D3A4DEE5FE7C5CCEB597B5691134647829AD30CA1E36EDBEAC32B6")
```

### Sign & Send transaction

```go
cosmosAccount = // create a cosmos account
toAddress = "cosmos1lkw6n8efpj7mk29yvajp......"
amount = "100000"
gasPrice = "0.01"
gasLimit = "80000"

chain = cosmos.NewChainWithRpc(rpcUrl, restUrl)
token = chain.DenomToken(prefix, denom)

signedTx, err = token.BuildTransferTx(account.PrivateKeyHex(), toAddress, gasPrice, gasLimit, amount)

txHash, err = chain.SendRawTransaction(signedTx)

```

--------------------------------------------------------------------------------
--------------------------------------------------------------------------------



## Ethereum Sign Transaction

```go
ethereumChain, err = eth.NewChainWithRpc(rpcUrl)

// make an transaction object
transaction = NewTransaction(nonce, gasPrice, gasLimit, to, value, data)
// transaction.MaxPriorityFeePerGas = "10000" // if send EIP1559 tx

// sign with hexed privatekey
signedTxObj, err = ethereumChain.SignTransaction(privateKeyHex, transaction)
// sign with ethereum account
signedTxObj, err = ethereumChain.SignTransactionWithAccount(ethAccount, transaction) 

signedHashString = signedTxObj.Value // signed transaction hash string
```

--------------------------------------------------------------------------------
--------------------------------------------------------------------------------

### Sui Merge Coin

```golang
var owner = "0x123abc"
var coinType = "0x2::sui::SUI" // Default
var targetAmount = (10,000,000,000).toString() // 10 SUI

// make request
var mergeRequest = chain.BuildMergeCoinRequest(owner, coinType, targetAmount)
println(mergeRequest.WillBeAchieved)
println(mergeRequest.EstimateAmount)
println("merging coin count = ", mergeRequest.CoinsCount)

// preview merge result
var mergePreview = chain.BuildMergeCoinPreview(mergeRequest)
println(mergePreview.SimulateSuccess)
println(mergePreview.EstimateGasFee)
println(mergePreview.WillBeAchieved)
println(mergePreview.EstimateAmount)

// sign & send
var txn = mergePreview.Transaction
var account = ...
var signedTxn = txn.SignWithAccount(account)
var hash = chain.SendRawTransaction(signedTxn)
```

### Sui Split Coin

```golang
// build transaction
var transaction = chain.BuildSplitCoinTransaction(owner, coinType, targetAmount)

// estimate gas fee
var estimateFee = chain.EstimateGasFee(transaction)

// sign & send
......
```


--------------------------------------------------------------------------------
--------------------------------------------------------------------------------



ComingChat substrate wallet SDK

## Build Android && Ios

* make buildAllAndroid
* make buildAllIOS
