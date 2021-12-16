# wallet-SDK

ComingChat substrate wallet SDK

## Struct

### Tx

#### 创建Tx

```java
Tx tx = Tx.newTx(metadataString);
```

| 参数           | 类型   | 描述                   |
| -------------- | ------ | ---------------------- |
| metadataString | string | metadata的16进制string |

throw err



#### 创建Balance转账

```java
Transaction t = tx.newBalanceTransferTx(dest, genesisHashString, amount, nonce, specVersion, transVersion);
```

| 参数              | 类型   | 描述         | 获取方式 |
| ----------------- | ------ | ------------ | -------- |
| dest              | string | 对方address  | 输入     |
| genesisHashString | string | 创世哈希     | 接口获取 |
| amount            | long   | 金额         | 用户输入 |
| nonce             | long   | nonc         | 接口获取 |
| specVersion       | int    | specVersion  | 接口获取 |
| transVersion      | int    | transVersion | 接口获取 |

throw err



#### 创建ChainX转账

```java
Transaction t = tx.newChainXBalanceTransferTx(dest, genesisHashString, amount, nonce, specVersion, transVersion);
```

参数同上

throw err



#### 创建NFT转账

```java
Transaction t = tx.newComingNftTransferTx(dest, genesisHashString, cid, nonce, specVersion, transVersion);
```

| 参数              | 类型   | 描述         | 获取方式 |
| ----------------- | ------ | ------------ | -------- |
| dest              | string | 对方address  | 输入     |
| genesisHashString | string | 创世哈希     | 接口获取 |
| cid               | long   | 金额         | 用户输入 |
| nonce             | long   | nonc         | 接口获取 |
| specVersion       | int    | specVersion  | 接口获取 |
| transVersion      | int    | transVersion | 接口获取 |

throw err



#### 创建XBTC转账

```java
Transaction t = tx.newXAssetsTransferTx(dest, genesisHashString, amount, nonce, specVersion, transVersion);
```

参数同上

throw err



#### 创建mini门限转账

```java
Transaction t = tx.NewThreshold(thresholdPublicKey, destAddress, aggSignature, aggPublicKey, controlBlock, message, scriptHash, genesisHashString, transferAmount, nonce, blockNumber, specVersion, transVersion);
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
| genesisHashString  | string | genesisHash                  | 接口获取                  |
| transferAmount     | int64  | 转账金额                     | 用户输入                  |
| nonce              | int64  | nonce                        | 接口获取                  |
| blockNumber        | int32  | blockNumber                  | 与scriptHash同接口获取    |
| specVersion        | int32  | specVersion                  | 接口获取                  |
| transVersion       | int32  | transVersion                 | 接口获取                  |



### Wallet

**Wallet 工具类方法（静态方法）**

#### 获取助记词：

```java
String mnemonicString = Wallet.genMnemonic();
```

throw err	

#### 创建钱包:

```java
Wallet wallet = Wallet.newWallet(seedOrPhrase , network);
```

| 参数         | 类型   | 描述           | 获取方式 |
| ------------ | ------ | -------------- | -------- |
| seedOrPhrase | string | 助记词或者私钥 |          |
| network      | int    | ss58           |          |




throw err

#### 地址转公钥

```java
String publicKey = Wallet.addressToPublicKey(address);
```

| 参数    | 类型   | 描述 | 获取方式 |
| ------- | ------ | ---- | -------- |
| address | string | 地址 |          |

#### 公钥转地址

```java
String address = Wallet.publicKeyToAddress(publicKey, network);
```

| 参数      | 类型   | 描述             | 获取方式 |
| --------- | ------ | ---------------- | -------- |
| publicKey | string | 公钥16进制字符串 |          |
| network   | int    | ss58             |          |



**Wallet类方法**

#### 获取公钥：

```java
String publicKeyHex = wallet.getPublicKeyHex();
byte[] publicKey = wallet.getPublicKey();
```

#### 获取私钥：

```java
String privateKey = wallet.getPrivateKeyHex();
```

throw err

#### 签名:

```java
byte[] signature = wallet.sign(message);
```

throw err

#### 获取地址:

```java
String address = wallet.getAddress();
```

throw err

### Transaction

#### 内部成员变量

```java
SignatureData    byte[]
PublicKey        byte[]
```

#### 获取需要签名内容

```java
byte[] signData = t.getSignData();
```

throw err

#### 获取Tx

```java
String sendTx = t.getTx();
```

throw err

## 交易构造流程

1. 导入助记词或私钥，生成钱包对象

2. 调用[创建钱包](#创建钱包)获得钱包对象wallet用于签名

3. 调用接口获取构造参数(nonce, specVersion, transVersion )

4. Metadata的获取根据接口判断

5. 调用[创建Tx](#创建Tx)获得tx对象用于构造交易

6. 通过tx对象创建交易t

7. 用钱包Sign方法对交易t的内容([获取需要签名内容](#获取需要签名内容))签名

8. 将签名方公钥([获取公钥](#获取公钥))和签名内容放入交易t中

9. 交易t调用[获取Tx](#获取Tx)获得SendTx

```java
try {
  String mnemonicString = Wallet.genMnemonic();
	Wallet wallet = Wallet.newWallet(mnemonicString , 44);
  Tx tx = Tx.newTx(metadataString);
  String dest = "5UczqUVGsoQpZnBCZkDtxvLxJ42KnUfaGTzPkQmZeAAug4s9";
  String genesisHashString = "0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f";
  long amount = 100000;
  Transaction t = tx.newBalanceTransferTx(dest, genesisHashString, amount, 0, 115, 1);
  byte[] signedMessage = wallet.sign(t.getSignData());
  t.setSignatureData(signedMessage);
  t.setPublicKey(wallet.getPublicKey());
  String sendTx = t.getTx();
}catch  (Exception e){
  
}
```

