# 红包合约 SDK

最新版的红包 SDK 代码设计：
- 基础的对象，接口 定义在 base 包
  - base.RedPacketAction  create/open/close 接口定义
  - base.RedPacketContract  合约接口
  - base.RedPacketDetail
- 各自链的 redpacket 操作在各自链的包里实现 base.RedPacketContract 接口
  - aptos.aptosRedPacketContract，通过 aptos.NewRedPacketContract 获得对象实例
  - eth.ethRedPacketContract，通过 eth.NewRedPacketContract 获得对象实例

Aptos 链 create 红包 example:
```go
chain := aptos.NewChainWithRestUrl(testNetUrl)
account, err := aptos.AccountWithPrivateKey(os.Getenv("private"))
if err != nil {
    panic(err)
}
contract := aptos.NewRedPacketContract(chain, os.Getenv("red_packet"))
// Aptos 红包合约目前仅支持 aptos 原生币，所以第一个参数传空即可
action, err := base.NewRedPacketActionCreate("", 5, "100000")
if err != nil {
    panic(err)
}
txHash, err := contract.SendTransaction(account, action)
if err != nil {
    panic(err)
}
txDetail, err := chain.FetchTransactionDetail(txHash)
if err != nil {
    panic(err)
}
```

Eth 链旧版的红包 SDK 使用的是 `eth.RedPacketAction`，此类已弃用，未来不再维护。

Eth 链也可以使用新版红包 SDK，基本逻辑相同。
```go
chain := eth.NewChainWithRpc(os.Getenv("rpc"))
account, err := eth.NewAccountWithMnemonic(os.Getenv(""))
if err != nil {
    panic(err)
}
contract := eth.NewRedPacketContract(chain, os.Getenv("red_packet"))
action, err := base.NewRedPacketActionCreate("", 5, "100000")
if err != nil {
    panic(err)
}
txHash, err := contract.SendTransaction(account, action)
if err != nil {
    panic(err)
}
txDetail, err := chain.FetchTransactionDetail(txHash)
if err != nil {
    panic(err)
}
```

