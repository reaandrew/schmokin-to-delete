.PHONY: proto
proto:
	 cd server && protoc --go_out=plugins=grpc:. *.proto

.PHONY: install_linter
install_linter:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0

.PHONY: lint
lint:
	golangci-lint run
	 
