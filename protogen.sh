#!/bin/bash

mkdir -p ./grpc_ts
mkdir -p ./grpc_go

if [[ "$GOBIN" == "" ]]; then
  if [[ "$GOPATH" == "" ]]; then
    echo "Required env var GOPATH is not set; aborting with error; see the following documentation which can be invoked via the 'go help gopath' command."
    go help gopath
    exit -1
  fi

  echo "Optional env var GOBIN is not set; using default derived from GOPATH as: \"$GOPATH/bin\""
  export GOBIN="$GOPATH/bin"
fi

echo "Compiling protobuf definitions"
protoc \
  --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
  --plugin=protoc-gen-go=${GOBIN}/protoc-gen-go \
  -I ./proto \
  -I/usr/local/include -I. \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --js_out=import_style=commonjs,binary:./grpc_ts \
  --go_out=plugins=grpc:./grpc_go \
  --grpc-gateway_out=logtostderr=true:./grpc_go \
  --ts_out=service=true:./grpc_ts \
  ./proto/brewery.proto

# protoc -I/usr/local/include -I. \
#   -I$GOPATH/src \
#   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#   --go_out=plugins=grpc:./grpc_go \
#   --grpc-gateway_out=logtostderr=true:./grpc_go \
#   ./proto/brewery.proto
