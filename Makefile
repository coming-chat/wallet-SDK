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
pkgCosmos =  $(pkgCore)/cosmos
pkgDoge =  $(pkgCore)/doge
pkgMSCheck = ${pkgCore}/multi-signature-check

pkgAll = ${pkgBase} ${pkgWallet} ${pkgPolka} ${pkgBtc} ${pkgEth} ${pkgCosmos} ${pkgMSCheck} ${pkgDoge}

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
	@make buildAllSDKIOS && make buildAllSDKAndroid
	@cd ${outdir} && zip -ry ${iosZipName}.${v}.zip Wallet.xcframework
	@cd ${outdir} && mkdir ${andZipName}.${v} && mv -f wallet.aar wallet-sources.jar ${andZipName}.${v}
	@cd ${outdir} && zip -ry ${andZipName}.${v}.zip ${andZipName}.${v}
	@cd ${outdir} && open .

# 给打包机打包归档用的
# ❯ ssh -T coming@192.168.3.84 << EOF
# ❯ cdsdk_wallet
# ❯ make archiveSDK v=0.0.2
# ❯ EOF
archiveSDK:
	@git reset --hard && git clean -id && git checkout -f main && git pull --rebase
	@make packageAll ${v}
	@cd ${outdir} && rm -rf ${v} && mkdir ${v}
	@cd ${outdir} && mv ${andZipName}.${v}.zip ${iosZipName}.${v}.zip ${v}
	@cd ${outdir}/${v} && echo `git log -1` > info.md
	@cd ${outdir}/${v} && echo "\n归档时间: " >> info.md
	@cd ${outdir}/${v} && echo `date +%FT%T%z` >> info.md
	@say '归档完成'
