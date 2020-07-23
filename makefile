
# path to the project.
ROOT_DIR:=			$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
LIB_FSPATH:=		$(GOPATH)/src/github.com/ctessum/cityaq
GO_OS:=				$(shell go env GOOS)
GO_ARCH:=			$(shell go env GOARCH)


print:
	@echo
	@echo ROOT_DIR : 	$(ROOT_DIR)
	@echo LIB_FSPATH : 	$(LIB_FSPATH)
	@echo GO_OS : 		$(GO_OS)
	@echo GO_ARCH : 	$(GO_ARCH)
	@echo


### BUILD
# Path to protoc binary 
PROTOC_FSPATH=$(LIB_FSPATH)/lib-protoc/bin/
export PATH:=$(PROTOC_FSPATH):$(PATH)

build:
	@echo
	@echo -- Download Protoc for OS  --
	cd $(LIB_FSPATH) && go build ./internal/download
	cd $(LIB_FSPATH) && ./download

	@echo
	@echo -- Download protoc-gen-go --
	cd $(LIB_FSPATH) && go get -u github.com/golang/protobuf/protoc-gen-go@v1.4.2

	# Generate the gRPC client/server code. (Information at https://grpc.io/docs/quickstart/go.html)
	@echo
	@echo -- Generate the gRPC client/server code --
	cd $(LIB_FSPATH) && protoc cityaq.proto --go_out=plugins=grpc:cityaqrpc
	cd $(LIB_FSPATH) && go build ./internal/addtags
	cd $(LIB_FSPATH) && ./addtags -file=cityaqrpc/cityaq.pb.go -tags=!js

	# Generate the gRPC WASM client/server code. (Information at https://grpc.io/docs/quickstart/go.html)
	@echo
	@echo -- Generate the gRPC WASM client/server code --
	# replace protoc-gen-go with the WASM version.
	cd $(LIB_FSPATH) && go get -u github.com/johanbrandhorst/grpc-wasm/protoc-gen-wasm@v0.0.0-20180613181153-d79a93c3901e
	cd $(LIB_FSPATH) && mv $(GOPATH)/bin/protoc-gen-wasm $(GOPATH)/bin/protoc-gen-go

	cd $(LIB_FSPATH) && protoc cityaq.proto --go_out=plugins=grpc:cityaqrpc
	cd $(LIB_FSPATH) && ./addtags -file=cityaqrpc/cityaq.wasm.pb.go -tags=js
	cd $(LIB_FSPATH) && rm addtags

	@echo
	@echo -- Client dep --
	cd $(LIB_FSPATH) && go get github.com/golang/mock/gomock 
	cd $(LIB_FSPATH) && go install github.com/golang/mock/mockgen

	@echo
	@echo -- Client WASM build --
	cd $(LIB_FSPATH) && GOOS=js GOARCH=wasm go build -o ./gui/html/cityaq.wasm ./gui/cmd/main.go
	cd $(LIB_FSPATH)/gui/html && ls -al

	@echo
	@echo -- Client Compression ( takes a long time ) --
	cd $(LIB_FSPATH) && go run internal/compress/main.go
	cd $(LIB_FSPATH) && rm ./gui/html/cityaq.wasm
	cd $(LIB_FSPATH)/gui/html && ls -al

	@echo
	@echo -- Client Update WASM runners --
	@echo  GOROOT:  $(GOROOT)
	cp $(GOROOT)/misc/wasm/wasm_exec.html $(LIB_FSPATH)/gui/html/index.html
	cp $(GOROOT)/misc/wasm/wasm_exec.js $(LIB_FSPATH)/gui/html/wasm_exec.js
	cd $(LIB_FSPATH)/gui/html && ls -al
	
	@echo
	@echo -- Client Pack into bindata --
	go-bindata --pkg cityaq -o $(LIB_FSPATH)/assets.go $(LIB_FSPATH)/gui/html/

gen:
	# The other way to build. Not working on Mac..
	cd $(LIB_FSPATH) && go generate ./...
	
server-run:
	cd $(LIB_FSPATH) && go run ./cmd .
	# https://127.0.0.1:1000/


	

	

	

