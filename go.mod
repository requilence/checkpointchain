module github.com/requilence/checkpointchain

go 1.12

require (
	github.com/allegro/bigcache v1.2.0 // indirect
	github.com/aristanetworks/goarista v0.0.0-20190607111240-52c2a7864a08 // indirect
	github.com/cosmos/cosmos-sdk v0.28.2-0.20190615115749-95b946b4559a
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.8.27
	github.com/gorilla/mux v1.7.0
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/prometheus/procfs v0.0.0-20190328153300-af7bedc223fb // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.0.3
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570 // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.31.5
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	google.golang.org/genproto v0.0.0-20190327125643-d831d65fe17d // indirect
	google.golang.org/grpc v1.19.1 // indirect
)

replace github.com/cosmos/cosmos-sdk => github.com/requilence/cosmos-sdk v0.28.2-0.20190616114425-c8da5418a9fe

replace github.com/ethereum/go-ethereum => github.com/etclabscore/ethereum.go-ethereum v0.0.0-20190521172606-ca7575d1f5d5
