#!/bin/sh

protoc \
-I=. ./proto/match/*.proto \
--js_out=import_style=commonjs:. \
--grpc-web_out=import_style=typescript,mode=grpcwebtext:.
