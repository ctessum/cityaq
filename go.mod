module github.com/ctessum/cityaq

go 1.13

require (
	github.com/andybalholm/brotli v0.0.0-20190821151343-b60f0d972eeb
	github.com/cenkalti/backoff v2.0.0+incompatible
	github.com/ctessum/geom v0.2.10-0.20200417141930-c1ad83ff7e0d
	github.com/ctessum/requestcache v1.0.1
	github.com/ctessum/requestcache/v4 v4.0.0
	github.com/ctessum/sparse v0.0.0-20181201011727-57d6234a2c9d
	github.com/ctessum/unit v0.0.0-20160621200450-755774ac2fcb
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/golang/gddo v0.0.0-20190904175337-72a348e765d2 // indirect
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/gonum/floats v0.0.0-20181209220543-c233463c7e82
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/improbable-eng/grpc-web v0.11.0
	github.com/johanbrandhorst/grpc-wasm v0.0.0-20180613181153-d79a93c3901e
	github.com/lpar/gzipped v1.1.0
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/paulmach/orb v0.1.5
	github.com/rs/cors v1.7.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spatialmodel/inmap v1.7.1-0.20200829195015-9143eacc1ccb
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	gonum.org/v1/gonum v0.0.0-20191009222026-5d5638e6749a
	gonum.org/v1/plot v0.0.0-20190615073203-9aa86143727f
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.25.0
	k8s.io/api v0.19.0
	k8s.io/client-go v0.19.0
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
