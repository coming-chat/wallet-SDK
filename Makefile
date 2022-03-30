VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/coming-chat/wallet-SDK

pkgCore = ${module}/core/eth

pkgEth =  $(pkgCore)

pkgPolka = ${module}/wallet

pkgBtc = ${module}/core/btc

pkgMSCheck = ${module}/core/multi-signature-check

buildAllSDKAndroid:
	gomobile bind -ldflags "-s -w" -target=android/arm,android/arm64 -o=${outdir}/wallet.aar ${pkgEth} ${pkgPolka} ${pkgBtc} ${pkgMSCheck}

buildAllSDKIOS:
	gomobile bind -ldflags "-s -w" -target=ios  -o=${outdir}/Wallet.xcframework ${pkgEth} ${pkgPolka} ${pkgBtc} ${pkgMSCheck}

packageAll:
	rm -rf ${outdir}/*
	@make buildAllAndroid && make buildAllIOS
	@cd ${outdir} && mkdir android && mv wallet* android
	@cd ${outdir} && tar czvf android.tar.gz android/*
	@cd ${outdir} && tar czvf Wallet.xcframework Wallet.xcframework/*