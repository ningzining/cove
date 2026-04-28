PROJECT:=cove-system

# make build: 编译二进制文件
.PHONY: build
build:
	export CGO_ENABLED=0 && \
	export GOOS=linux && \
	export GOARCH=amd64 && \
	go build -ldflags="-w -s" -o $(PROJECT) ./cmd/cove-system

# make swag: 生成文档
.PHONY: swag
swag:
	cd ./internal/system && \
	swag init --pd -o ./docs/ --instanceName system -g ./system.go

# make run: 本地运行程序
.PHONY: run
run:
	go run cmd/cove-system/main.go -f etc/system.local.yaml

# make help: 查看命令帮助文档
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^#//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'	