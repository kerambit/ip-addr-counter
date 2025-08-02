.PHONY: build-mac clean

build-mac:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/counter-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/counter-amd64 .
	lipo -create -output bin/counter bin/counter-arm64 bin/counter-amd64
	rm bin/counter-arm64
	rm bin/counter-amd64

clean:
	rm -rf bin/