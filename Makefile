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
pkgSolana =  $(pkgCore)/solana
pkgAptos =  $(pkgCore)/aptos
pkgSui =  $(pkgCore)/sui
pkgStarcoin =  $(pkgCore)/starcoin
pkgMSCheck = ${pkgCore}/multi-signature-check

pkgAll = ${pkgBase} ${pkgWallet} ${pkgPolka} ${pkgBtc} ${pkgEth} ${pkgCosmos} ${pkgMSCheck} ${pkgDoge} ${pkgSolana} ${pkgAptos} ${pkgSui} ${pkgStarcoin}

iosPackageName=Wallet.xcframework

buildAllSDKAndroid:
	gomobile bind -ldflags "-s -w" -target=android/arm,android/arm64 -o=${outdir}/wallet.aar ${pkgAll}

buildAllSDKIOS:
	GOOS=ios gomobile bind -ldflags "-s -w" -v -target=ios/arm64  -o=${outdir}/${iosPackageName} ${pkgAll}

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

iosReposity=${outdir}/Wallet-iOS
iosCopySdk:
	@cd ${iosReposity} \
		&& rm -rf Sources/* \
		&& cp -Rf ../${iosPackageName} Sources \
		&& cp -Rf ../${iosPackageName}/ios-arm64/Wallet.framework/Versions/A/Headers Sources

iosPublishVersion:
ifndef v
	@echo 发布 iOS 包需要指定一个版本，例如 make publishIOSVersion v=1.0.1
	@exit 1
endif
	@make iosCopySdk
	@cd ${iosReposity} \
		&& git add --all \
		&& git commit -m 'Auto Publish ${v}' -m "refer to `git rev-parse HEAD`" \
		&& git tag -f ${v} \
		&& git push origin main tag ${v} --force
	@make iosPublishMain

iosPublishMain:
	@make iosCopySdk
	@cd ${iosReposity} \
		&& rm -rf Sources/Wallet.xcframework/ios-arm64_x86_64-simulator \
		&& git add --all \
		&& git commit -m 'Auto Publish Develop SDK' \
		&& git push origin main
