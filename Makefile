VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/coming-chat/wallet-SDK

pkgCore = ${module}/core/eth

pkgAll =  $(pkgCore)
#$(pkgCore) $(pkgUtil) $(pkgGasNow) $(pkgConstants)

buildAllAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/ethwalletcore.aar ${pkgAll}
buildAllIOS:
	gomobile bind -ldflags "-s -w" -target=ios  -o=${outdir}/ethwalletcore.xcframework ${pkgAll}

packageAll:
	rm -rf ${outdir}/*
	@make buildAllAndroid && make buildAllIOS
	@cd ${outdir} && mkdir android && mv eth-wallet* android
	@cd ${outdir} && tar czvf android.tar.gz android/*
	@cd ${outdir} && tar czvf eth-wallet.xcframework eth-wallet.xcframework/*