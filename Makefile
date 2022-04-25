VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/coming-chat/wallet-SDK

pkgCore = ${module}/core
pkgBase = ${pkgCore}/base
pkgWallet = ${pkgCore}/wallet
pkgPolka = ${pkgCore}/polka
pkgBtc = ${pkgCore}/btc
pkgEth =  $(pkgCore)/eth
pkgMSCheck = ${pkgCore}/multi-signature-check

pkgAll = ${pkgBase} ${pkgWallet} ${pkgPolka} ${pkgBtc} ${pkgEth} ${pkgMSCheck} 

buildAllSDKAndroid:
	gomobile bind -ldflags "-s -w" -target=android/arm,android/arm64 -o=${outdir}/wallet.aar ${pkgAll}

buildAllSDKIOS:
	GOOS=ios gomobile bind -ldflags "-s -w" -target=ios  -o=${outdir}/Wallet.xcframework ${pkgAll}

# 使用: make packageAll v=1.4
# 结果: out 目录下将产生两个压缩包 wallet-SDK-ios.1.4.zip 和 wallet-SDK-android.1.4.zip 
iosZipName=wallet-SDK-ios
andZipName=wallet-SDK-android
packageAll:
	# rm -rf ${outdir}/*
	@cd ${outdir} && rm -f wallet-SDK-*.zip && rm -rf ${andZipName}.${v}
	@make buildAllSDKAndroid && make buildAllSDKIOS
	@cd ${outdir} && zip -ry ${iosZipName}.${v}.zip Wallet.xcframework
	@cd ${outdir} && mkdir ${andZipName}.${v} && mv -f wallet.aar wallet-sources.jar ${andZipName}.${v}
	@cd ${outdir} && zip -ry ${andZipName}.${v}.zip ${andZipName}.${v}
	@cd ${outdir} && open .
