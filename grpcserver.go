package cityaq

import (
	"net/http"
	"strings"

	"github.com/ctessum/cityaq/cityaqrpc"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCSServer is a GRPC server for city AQ simulation data.
type GRPCServer struct {
	CityAQ
	grpcServer   *grpcweb.WrappedGrpcServer
	staticServer http.Handler

	Log logrus.FieldLogger
}

// NewGRPCServer creates a new GRPC server for c.
func NewGRPCServer(c *CityAQ, prefix string) *GRPCServer {
	gs := grpc.NewServer()
	cityaqrpc.RegisterCityAQServer(gs, c)
	s := new(GRPCServer)
	s.grpcServer = grpcweb.WrapServer(gs)
	s.staticServer = wasmContentTypeSetter(http.FileServer(
		&assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    "gui/html",
		},
	))
	return s
}

func wasmContentTypeSetter(fn http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.URL.Path, ".wasm") {
			w.Header().Set("content-type", "application/wasm")
		}
		fn.ServeHTTP(w, req)
	}
}

func (s *GRPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Redirect gRPC and gRPC-Web requests to the gRPC-Web Websocket Proxy server
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		if s.Log != nil {
			s.Log.WithFields(logrus.Fields{
				"url":  r.URL.String(),
				"addr": r.RemoteAddr,
			}).Info("cityaq grpc request")
		}
		s.grpcServer.ServeHTTP(w, r)
		return
	} else {
		if s.Log != nil {
			s.Log.WithFields(logrus.Fields{
				"url":  r.URL.String(),
				"addr": r.RemoteAddr,
			}).Info("cityaq static request")
		}
		s.staticServer.ServeHTTP(w, r)
	}
}
