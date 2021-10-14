module github.com/ctessum/cityaq

go 1.13

require (
	github.com/andybalholm/brotli v0.0.0-20190821151343-b60f0d972eeb
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/coreos/bbolt v1.3.1-coreos.6 // indirect
	github.com/coreos/etcd v3.3.10+incompatible // indirect
	github.com/ctessum/geom v0.2.10
	github.com/ctessum/requestcache/v4 v4.0.0
	github.com/ctessum/sparse v0.0.0-20181201011727-57d6234a2c9d
	github.com/ctessum/unit v0.0.0-20160621200450-755774ac2fcb
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/go-ini/ini v1.39.0 // indirect
	github.com/golang/gddo v0.0.0-20190904175337-72a348e765d2 // indirect
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/gonum/floats v0.0.0-20181209220543-c233463c7e82
	github.com/improbable-eng/grpc-web v0.11.0
	github.com/johanbrandhorst/grpc-wasm v0.0.0-20180613181153-d79a93c3901e
	github.com/lpar/gzipped v1.1.0
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/paulmach/orb v0.1.6
	github.com/rs/cors v1.7.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spatialmodel/inmap v1.9.0
	golang.org/x/net v0.0.0-20210505214959-0714010a04ed
	gonum.org/v1/gonum v0.0.0-20191009222026-5d5638e6749a
	gonum.org/v1/plot v0.0.0-20190615073203-9aa86143727f
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/ini.v1 v1.39.0 // indirect
	gopkg.in/pipe.v2 v2.0.0-20140414041502-3c2ca4d52544 // indirect
	k8s.io/api v0.20.1
	k8s.io/client-go v0.20.1
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
