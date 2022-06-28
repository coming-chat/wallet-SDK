package ignoredcopy

type AccountCase struct {
	Mnemonic string
	Address  string
}

type AccountGroup struct {
	BtcMainnet  AccountCase
	BtcSignet   AccountCase
	Cosmos      AccountCase
	Terra       AccountCase
	DogeMainnet AccountCase
	DogeTestnet AccountCase
	Ethereum    AccountCase
	Polka0      AccountCase
	Polka2      AccountCase
	Polka44     AccountCase
}

// 这些助记词是完全公开的，所以不要把你资产存入这里面
// These mnemonics are completely public, so don't deposit your assets here
const m1 = "unaware oxygen allow method allow property predict various slice travel please priority"
const mterra = "canyon young easy visa antenna address zone maple captain garden faith crawl tomorrow left risk identify impose miss baby whale nest assume clap trial"

var Accounts = AccountGroup{
	BtcMainnet:  AccountCase{m1, "bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sg"},
	BtcSignet:   AccountCase{m1, "tb1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ystm5728"},
	Cosmos:      AccountCase{m1, "cosmos19jwusy7lm8v5kqay8qjml79hs6e30t8j7ygm8r"},
	Terra:       AccountCase{mterra, "terra1swy7k7r0jv4rmyjslp35pf0dfp0cs92c8mdwlr"},
	DogeMainnet: AccountCase{m1, "DJhF8ahvTfGhqcLEn7sN4gJMJVVbmfwxkU"},
	DogeTestnet: AccountCase{m1, "nhkJrbSqPdjRiauRowWpK5teYMstkMp4M6"},
	Ethereum:    AccountCase{m1, "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2"},
	Polka0:      AccountCase{m1, "12jrfZLTddDxRQAjoSkWurDyEPxPdkhPcgU2AGxFHbgBpyHZ"},
	Polka2:      AccountCase{m1, "EKBBYRGQCyQjWyfcWWZfekpXNEyk7xRzZaHPeErDJsAPeiD"},
	Polka44:     AccountCase{m1, "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH"},
}
