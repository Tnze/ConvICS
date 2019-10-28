#/bin/bash
export GOARCH=wasm
export GOOS=js
go build -o docs/wasm/xatu.wasm ./xatu