PROJECT:=cove-sys

# make build: 编译二进制文件
.PHONY: build
build: swag
	export CGO_ENABLED=0 && \
	export GOOS=linux && \
	export GOARCH=amd64 && \
	go build -ldflags="-w -s" -o $(PROJECT) .

# make swag: 生成文档
.PHONY: swag
swag:
	swag init --pd -o ./docs/admin/ --instanceName $(PROJECT)

# make run: 本地运行程序
.PHONY: run
run:
	go run main.go -f etc/sys.local.yaml

# make help: 查看命令帮助文档
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^#//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'	