VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/coming-chat/wallet-SDK

pkgCore = ${module}/core

pkgEth =  $(pkgCore)/eth

pkgPolka = ${pkgCore}/wallet

pkgBtc = ${pkgCore}/btc

pkgMSCheck = ${pkgCore}/multi-signature-check

pkgAll = ${pkgEth} ${pkgPolka} ${pkgBtc} ${pkgMSCheck}

buildAllSDKAndroid:
	gomobile bind -ldflags "-s -w" -target=android/arm,android/arm64 -o=${outdir}/wallet.aar ${pkgAll}

buildAllSDKIOS:
	gomobile bind -ldflags "-s -w" -target=ios  -o=${outdir}/Wallet.xcframework ${pkgAll}

# 使用: make packageAll v=1.4
# 结果: out 目录下将产生两个压缩包 wallet-SDK-ios.1.4.zip 和 wallet-SDK-android.1.4.zip 
iosZipName=wallet-SDK-ios
andZipName=wallet-SDK-android
packageAll:
	# rm -rf ${outdir}/*
	@make buildAllSDKAndroid && make buildAllSDKIOS
	@cd ${outdir} && zip -ry ${iosZipName}.${v}.zip Wallet.xcframework
	@cd ${outdir} && mkdir ${andZipName}.${v} && mv -f wallet.aar wallet-sources.jar ${andZipName}.${v}
	@cd ${outdir} && zip -ry ${andZipName}.${v}.zip ${andZipName}.${v}
	@cd ${outdir} && open .