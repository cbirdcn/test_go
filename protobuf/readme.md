# protobuf与go

## 笔记

注意事项已在代码中标记，其他可见下面参考内容，这里只记载一些操作

下载protoc二进制（注意不要下载到arm版本，文件名很像）：[protoc-3.11.2-linux-x86_64](https://github.com/protocolbuffers/protobuf/releases/download/v3.11.2/protoc-3.11.2-linux-x86_64.zip)

直接解压到`/usr/local`。因为一般情况下,`$GOPATH`都包含了`/usr/local`，解压的文件包含了`bin`和`include`等文件夹会直接合并到`/usr/local`的这些文件夹中。这样就避免解压到独立文件夹再加环境变量了。

然后安装`protoc-gen-go`，它可以将 `.proto` 文件转换为 Golang 代码。

根据不同的Go版本，`1.16+`使用`go install github.com/golang/protobuf/protoc-gen-go`，否则使用`go get -u github.com/golang/protobuf/protoc-gen-go`。`protoc-gen-go`将自动安装到 `$GOPATH/bin` 目录下

编写好多个proto文件，注意option用于确定pb文件生成位置和生成时的package名。

编写proto时，注意enum、string、bool的值为零值（0、""、false）时，不会进行编解码。这是protobuf的设计。这些不会编解码的数据仍然可以被读取，因为他们直接使用了默认值。

生成pb.go代码到pb文件夹内，生成代码的package名为pb：
```shell
cd ./pb
protoc --go_out=. *.proto
```
还可以在参数中指定多个源文件、目标路径等

多种类型数据的编解码过程可以看marshal文件夹。

rpc可以看rpc_server和rpc_client，结合main.proto。注意rpc_server中可被远程客户端访问方法有诸多要求。

这里不做支持http的行为，具体的可以看Kratos。

## 参考

[Go Protobuf 简明教程](https://geektutu.com/post/quick-go-protobuf.html)

[Go Generated Code Guide](https://protobuf.dev/reference/go/go-generated/#package)

[Language Guide (proto 3)](https://protobuf.dev/programming-guides/proto3/)

[Protocol Buffer Basics: Go](https://protobuf.dev/getting-started/gotutorial/)

[Define dictionary in protocol buffer](https://stackoverflow.com/questions/11474416/define-dictionary-in-protocol-buffer/11486640#11486640)

[Language Guide (proto 2)](https://protobuf.dev/programming-guides/proto2/#extensions)

[Go Rpc Examples](https://github.com/grpc/grpc-go/tree/master/examples)

[官方提供的proto包常规用法](https://github.com/golang/protobuf/tree/master/proto)

[官方提供的特殊类型ptypes包用法](https://github.com/golang/protobuf/tree/master/ptypes)

[Kratos的beershop项目](https://github.com/go-kratos/beer-shop/tree/main)

以及一些小问题处理

[gob: type not registered for interface](https://stackoverflow.com/questions/21934730/gob-type-not-registered-for-interface-mapstringinterface)

[proto3特殊默认值编解码被隐藏](https://blog.csdn.net/cs10239dn/article/details/125166742)