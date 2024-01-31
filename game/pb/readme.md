# protobuf与rpc

下载protoc二进制（注意不要下载到arm版本，文件名很像）：[protoc-3.11.2-linux-x86_64](https://github.com/protocolbuffers/protobuf/releases/download/v3.11.2/protoc-3.11.2-linux-x86_64.zip)

直接解压到 `/usr/local`。因为一般情况下,`$GOPATH`都包含了 `/usr/local`，解压的文件包含了 `bin`和 `include`等文件夹会直接合并到 `/usr/local`的这些文件夹中。这样就避免解压到独立文件夹再加环境变量了。

然后安装 `protoc-gen-go`，它可以将 `.proto` 文件转换为 Golang 代码。

根据不同的Go版本，`1.16+`使用 `go install github.com/golang/protobuf/protoc-gen-go`，否则使用 `go get -u github.com/golang/protobuf/protoc-gen-go`。`protoc-gen-go`将自动安装到 `$GOPATH/bin` 目录下

编写好多个proto文件，注意option用于确定pb文件生成位置和生成时的package名。

生成pb.go代码到pb文件夹内，生成代码的package名为pb

```shell
cd ./pb
protoc --go_out=. *.proto
```

status_code是错误码
