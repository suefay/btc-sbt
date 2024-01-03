#!/usr/bin/make -f

export GOPROXY = https://goproxy.io

build: go.sum
	@echo "build btc-sbt"
	@go build -mod=readonly -o $${GOBIN-$${GOPATH-$$HOME/go}/bin}/btc-sbt

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
