pkgs          = $(shell go list ./... | grep -v /tests | grep -v /vendor/ | grep -v /common/)

test:
	@echo " >> running tests"
	@go test  -cover $(pkgs)

genproto:
	protoc proto/*/*.proto --go_out=plugins=grpc:.

genmock:
	go generate ./...

download:
	go mod download

tidy:
	go mod tidy
