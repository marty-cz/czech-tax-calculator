GOCMD=go
BINARY_NAME=czech-tax-calculator
VERSION?=0.0.0

tidy:
	$(GOCMD) mod tidy

test: 
#	go test -v ./... -short
	go test -v ./...

build: tidy test
	mkdir -p out/bin
#	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o out/bin/${BINARY_NAME}-darwin cmd/${BINARY_NAME}/main.go
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o out/bin/${BINARY_NAME}-linux cmd/${BINARY_NAME}/main.go
#	CGO_ENABLED=0 GOARCH=amd64 GOOS=window go build -o out/bin/${BINARY_NAME}-windows cmd/${BINARY_NAME}/main.go


buildAndRun: build
	./out/bin/${BINARY_NAME}-linux --stock-input ./examples/Ucetni-kniha-Akcie.xlsx --crypto-input ./examples/Ucetni-kniha-Crypto.xlsx