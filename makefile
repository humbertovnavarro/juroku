build:
	GOOS=js GOARCH=wasm go build -buildvcs=false -o public/main.wasm
