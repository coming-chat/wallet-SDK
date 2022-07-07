# wallet-SDK

|                       | Bitcoin | Ethereum | Polka | Cosmos | Doge | Solana |
| --------------------- | ------- | -------- | ----- | ------ | ------ | ------ |
| import mnemonic       | ✅     | ✅      | ✅   |✅   |✅   |✅   |
| import keystore       | ❌     | ❌      | ✅  |❌   |❌   |❌   |
| pri/pub key & address | ✅     | ✅      | ✅   |✅   |✅   |✅   |
| multi network         | ✅     | ✅      | ✅   |✅   |❌   |✅   |
| publicKey to address  | ✅     | ✅      | ✅   |✅   |✅   |✅   |
| address to publicKey  | ❌     | ❌      | ✅   |❌   |❌   |✅   |
| sign data             | ☑️    | ☑️     | ✅   |☑️   |☑️   |☑️   |
|                       |         |          |       |        |      |        |
| query balance | ✅ | ✅ | ✅ |✅   |✅   |✅   |
| fetch transaction detail | ✅ | ✅ | ✅ |✅   |✅   |✅   |
| gas fee | ❌ | ✅ | ✅ |☑️   |✅   |✅   |
| send taw transaction | ✅ | ✅ ☑️ | ✅ ☑️ |✅   |✅   |✅   |
| multi token | ❌ | ✅ erc20 | ✅ XBTC |✅   |❌   |❌   |

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



ComingChat substrate wallet SDK

## build android && ios

* make buildAllAndroid
* make buildAllIOS

## Struct

### Wallet

#### **Wallet 工具类方法（静态方法）**

* ##### 获取助记词：

  ```java
  String mnemonicString = Wallet.genMnemonic();
  ```

  throw err	



* ##### 创建钱包:

  1. 从助记词或私钥

    ```java
    Wallet wallet = Wallet.newWallet(seedOrPhrase);
    ```

    | 参数         | 类型   | 描述           | 获取方式 |
    | ------------ | ------ | -------------- | -------- |
    | seedOrPhrase | string | 助记词或者私钥 |          |
    | network      | int    | ss58           |          |

    throw err

  

  2. 从keystore
  
    ```java
    Wallet wallet = Wallet.NewWalletFromKeyStore(keyStoreJson, password);
    ```
  
    | 参数         | 类型   | 描述     | 获取方式 |
    | ------------ | ------ | -------- | -------- |
    | keyStoreJson | string | keystore |          |
    | password     | string | 密码     |          |
  
    throw err



* ##### 地址转公钥

  ```java
  String publicKey = Wallet.addressToPublicKey(address);
  ```

  | 参数    | 类型   | 描述 | 获取方式 |
  | ------- | ------ | ---- | -------- |
  | address | string | 地址 |          |

​		throw err

* ##### 公钥转地址

    ```java
    String address = Wallet.publicKeyToAddress(publicKey, network);
    ```

    | 参数      | 类型   | 描述             | 获取方式 |
    | --------- | ------ | ---------------- | -------- |
    | publicKey | string | 公钥16进制字符串 |          |
    | network   | int    | ss58             |          |

​		throw err

#### **Wallet类方法**

  * ##### 获取公钥：

    ```java
    String publicKeyHex = wallet.getPublicKeyHex();
    byte[] publicKey = wallet.getPublicKey();
    ```


  * ##### 获取私钥：

    ```java
    String privateKey = wallet.getPrivateKeyHex();
    ```

    throw err



  * ##### 签名:

    ```java
    byte[] signature = wallet.sign(message, password);
    ```

    | 参数     | 类型   | 描述                                            | 获取方式 |
    | -------- | ------ | ----------------------------------------------- | -------- |
    | message  | []byte | 签名内容                                        |          |
    | password | string | keystore 密码，如果是助记词或者私钥创建则随便填 |          |

    throw err
    
    
    
  * ##### 检查Keystore密码

    ```java
    bool isCorrect = wallet.checkPassword(password);
    ```

    | 参数     | 类型   | 描述 | 获取方式 |
    | -------- | ------ | ---- | -------- |
    | password | string |      |          |
    

throw err



  * ##### 获取地址:

    ```java
    String address = wallet.getAddress(network);
    ```

    | 参数    | 类型 | 描述 | 获取方式 |
    | ------- | ---- | ---- | -------- |
    | network | int  | ss58 |          |

    throw err



### Tx

  * ##### 创建Tx

    ```java
    Tx tx = Tx.NewTx(metadataString);
    ```

    | 参数           | 类型   | 描述                   | 获取方式 |
    | -------------- | ------ | ---------------------- | -------- |
    | metadataString | string | metadata的16进制string | 接口获取 |

    throw err
    
  * ##### 从Hex创建Transaction

    ```
    Transaction t = tx.newTransactionFromHex(txHex)
    ```

    | 参数     | 类型   | 描述           | 获取方式             |
    | -------- | ------ | -------------- | -------------------- |
    | txHex    | string | tx的16进制数据 | 从Wallet Connect获取 |

    throw err



  * ##### 创建Balance转账

    ```java
    Transaction t = tx.newBalanceTransferTx(dest, amount);
    ```

    | 参数   | 类型   | 描述        | 获取方式 |
    | ------ | ------ | ----------- | -------- |
    | dest   | string | 对方address | 输入     |
    | amount | string | 金额        | 用户输入 |
    
    throw err



  * ##### 创建ChainX转账

    ```java
    Transaction t = tx.newBalanceTransferTx(dest, amount);
    ```

    参数同上

    throw err



  * ##### 创建NFT转账

    ```java
    Transaction t = tx.newComingNftTransferTx(dest, cid);
    ```

    | 参数 | 类型   | 描述        | 获取方式 |
    | ---- | ------ | ----------- | -------- |
    | dest | string | 对方address | 输入     |
    | cid  | long   | cid         | 用户输入 |
    
    throw err



  * ##### 创建XBTC转账

    ```java
    Transaction t = tx.newXAssetsTransferTx(dest, amount);
    ```

    参数同上

    throw err



  * ##### 创建mini门限转账

    ```java
    Transaction t = tx.NewThreshold(thresholdPublicKey, destAddress, aggSignature, aggPublicKey, controlBlock, message, scriptHash, transferAmount, blockNumber);
    ```

    | 参数               | 类型   | 描述                         | 获取方式                  |
    | :----------------- | ------ | ---------------------------- | ------------------------- |
    | thresholdPublicKey | string | 门限钱包公钥                 | generateThresholdPubkey() |
    | destAddress        | string | 对方地址                     |                           |
    | aggSignature       | string | 聚合签名                     | getAggSignature()         |
    | aggPublicKey       | string | 聚合公钥                     | getAggPubkey()            |
    | controlBlock       | string | controlBlock                 | generateControlBlock()    |
    | message            | string | 576520617265206c6567696f6e21 | 写死                      |
    | scriptHash         | string | scriptHash                   | 接口获取                  |
    | transferAmount     | string | 转账金额                     | 用户输入                  |
    | blockNumber        | int32  | blockNumber                  | 与scriptHash同接口获取    |
    
    throw err





### Transaction

  * ##### 获取未签名Tx

    ```java
    String unSignTx = t.getUnSignTx();
    ```
    throw err
    此方法用于预估矿工费使用



  * ##### 获取需要签名内容

    ```java
    byte[] signData = t.getSignData(genesisHashString, nonce, specVersion, transVersion);
    ```

    | 参数              | 类型   | 描述        | 获取方式 |
    | ----------------- | ------ | ----------- | -------- |
    | genesisHashString | string | genesisHash | 接口获取 |
    | nonce             | int64  | nonce       | 接口获取 |
    | specVersion       | int32  |             | 接口获取 |
    | transVersion      | int32  |             | 接口获取 |

    throw err

  * ##### 获取Tx

    ```java
    String sendTx = t.getTx(signerPublicKey, signatureData);
    ```

    | 参数            | 类型   | 描述       | 获取方式                                |
    | --------------- | ------ | ---------- | --------------------------------------- |
    | signerPublicKey | byte[] | 签名者公钥 | [获取公钥](#获取公钥)的getPublicKey方法 |
    | signatureData   | byte[] | 签名结果   |                                         |

    throw err





## 交易构造流程

1. 导入助记词或私钥，生成钱包对象

2. 调用[创建钱包](#创建钱包)获得钱包对象wallet用于签名

3. 调用接口获取构造参数(nonce, specVersion, transVersion )

4. Metadata的获取根据接口判断

5. 调用[创建Tx](#创建Tx)获得tx对象用于构造交易

6. 通过tx对象创建交易t

7. 用钱包Sign方法对交易t的内容([获取需要签名内容](#获取需要签名内容))签名

8. 将签名方公钥([获取公钥](#获取公钥))和签名内容放入GetTx的方法中

9. 交易t调用[获取Tx](#获取Tx)获得SendTx

```java
try {
  String mnemonicString = Wallet.genMnemonic();
	Wallet wallet = Wallet.newWallet(mnemonicString);
  Tx tx = Tx.NewTx(metadataString);
  String dest = "5UczqUVGsoQpZnBCZkDtxvLxJ42KnUfaGTzPkQmZeAAug4s9";
  String genesisHashString = "0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f";
  long amount = "100000";
  Transaction t = tx.newBalanceTransferTx(dest, amount);
  byte[] signedMessage = wallet.sign(t.getSignData(genesisHashString, 0, 115, 1));
  String sendTx = t.getTx(wallet.getPublicKey(), signedMessage);
}catch  (Exception e){
  
}
```
### Eth 交易golang代码，两端调用参考下，具体看sdk打包后源码

#### goland 代码

```golang
package main

import (
	"fmt"
	"math/big"

	coin "github.com/0xhelloweb3/go-coin-wallet/eth"
	"github.com/ethereum/go-ethereum/common"
)

var (
	rpcUrl = "https://data-seed-prebsc-1-s1.binance.org:8545"

	walletAddress = "0xB553803EE21b486BB86f2A63Bd682529Aa7FCE8D"

	privateKey = ""

	// bsc 测试网 busd 合约地址
	busdContractAddress = "0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee"
)


func main() {
	wallet := coin.NewEthChain()
	wallet.CreateRemote(rpcUrl)

	// 获取主网代币 BNB 余额
	balance, _ := wallet.Balance(walletAddress)
	fmt.Printf("bnb balance: %v\n", balance)

	// 获取 Erc20代币 余额
	busdBalance, _ := wallet.TokenBalance(busdContractAddress, walletAddress)

	tokenDecimal, err := wallet.TokenDecimal(busdContractAddress, walletAddress)
	fmt.Printf("busdBalance balance: %v, decimal: %v \n", busdBalance, tokenDecimal)

	if err != nil {
		fmt.Printf("get busdt balance error: %v\n", err)
		return
	}
	nonce, _ := wallet.Nonce(walletAddress)

	// 构造多笔交易则nonce + 1
	callMethodOpts := &coin.CallMethodOpts{
		Nonce: nonce,
	}

	// erc20 代币转账
	buildTxResult, err := wallet.BuildCallMethodTx(
		privateKey,
		busdContractAddress,
		coin.Erc20AbiStr,
		// 调用的合约方法名
		"transfer",
		callMethodOpts,
		// 转账目标地址
		common.HexToAddress("0x4165FD787ffF9f659A3B9A239a1d70fc5B8aB6d1"),
		big.NewInt(10000000000))

	if err != nil {
		fmt.Printf("build call method tx error: %v\n", err)
	}
	// 发送交易
	sendTxResult, err := wallet.SendRawTransaction(buildTxResult.TxHex)
	if err != nil {
		fmt.Printf("send raw transaction error: %v\n", err)
	}
	// 打印交易hash
	fmt.Printf("sendTxResult: %v\n", sendTxResult)

	// 检测 transfer 事件， fromBlock 和 toBlock 参数可以不传
	eventlogs, _ := wallet.FindLogs(busdContractAddress, coin.Erc20AbiStr, "Transfer",
		big.NewInt(17174691), big.NewInt(17174691), nil)
	fmt.Printf("eventlogs: %v\n", eventlogs)
}
```

#### Android 调用代码 (ERC20 转账)

```golang
 // 钱包地址
 private String walletAddress = "0xB553803EE21b486BB86f2A63Bd682529Aa7FCE8D";
 // 钱包对应的私钥
 private String privateKey = "";
 // busd 测试网合约地址
 private String busdContractAddress = "0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee";
 // bsc测试网rpc 节点
 bscRpc ="https://data-seed-prebsc-1-s1.binance.org:8545"

 EthChain ethChain = new EthChain();
 ethChain.createRemote(bscRpc);
 // 查询bnb余额
 long balance = ethChain.balance(walletAddress);
 Toast.makeText(getApplicationContext(),String.valueOf(balance),Toast.LENGTH_LONG).show();
 long Nonce =  ethChain.nonce(walletAddress);
 CallMethodOpts opts = new CallMethodOpts();
 // 设置nonce,可选
 opts.setNonce(Nonce);
 // 设置 gas 可选
 BuildTxResult buildTxResult = ethChain.buildCallMethodTx(
                    privateKey,
                    busdContractAddress,
                    Eth.getErc20AbiStr(),
                    "transfer",
                    opts,
                    {\"toAddress\":\"0x178a8AB44b71858b38Cc68f349A06f397A73bFf5\", \"amount\":1000000, \"method\":\"transfer\"}");
  Log.d("tx hex",buildTxResult.getTxHex());
  String tx = ethChain.sendRawTransaction(buildTxResult.getTxHex());
  Log.d("send tx success", tx);
```


#### PolkaChain Metadata 相关用法

```golang
// 只通过 rpcUrl 来初始化
chain = NewPolkaChain(rpc, scan)

// 通过传入本地缓存的 metadata 来初始化, 如果没有可以传空
chain = NewPolkaChainWithRpc(rpc, scan, metadata)

// 获取 metadata string
metadata, err = chain.GetMetadataString()

// 刷新 metadata string (可以从返回值取得最新的 metadata)
metadata, err = chain.ReloadMetadata()
```