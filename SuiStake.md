## Sui Stake Usage

### fetch validator list

```golang
var chain = sui.NewChainWithRpcUrl(rpcUrl)
var state = chain.GetValidatorState()
print(state.Epoch/TotalStaked/TotalRewards)
print("validator list = ", state.Validators)

// show validator infomation.
var validator = state.Validators[idx]
print(validator.Name/Address/ImageUrl/APY ....)
```



### fetch user staked delagation

```golang
var ownerAddress = "0x123456...."
var stakes = chain.GetDelegatedStakes(ownerAddress)

var stake = stakes[idx]
print(stake.Principal/StakeId/ValidatorAddress/Status ...)
if stake.Status == DelegationStatusPending {
  // pending
} else if stake.Status == DelegationStatusActived {
  print(stake.DelegationId/EarnedAmount ...)
}
```



### add stake delegation

```go
var validator = state.Validators[idx]

var amount = "1000000000" // 1 SUI
var txn = chain.AddDelegation(ownerAddress, amount, validator.Address)

var hash = // sign & send `txn`
```



### withdraw stake delegation

```golang
var stake = stakes[idx]

if stake.Status == DelegationStatusPending {
  return error("The pending stake delegation cannot be withdraw")
} else if stake.Status == DelegationStatusActived {
  var txn = chain.WithdrawDelegation(ownerAddress, stake.delegationId, stake.stakeId)
  
	var hash = // sign & send `txn`
}
```



