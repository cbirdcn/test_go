# db服务

## 简介

基于`go.mongodb.org/mongo-driver/mongo`提供对`mongodb`的访问。

`model`用于模型管理，`operation`用于api操作，每个`collection`需要自成一个文件。

`operation`中`UseCollection()`用于切换集合，每当外部（比如dbServer）想调用此包中的方法时，必须先调用切换集合的方法。

`model`中定义的数据模型，会在`operation`、`dbServer`中使用，比如读取DB数据并保存到本地struct变量中。

读写错误需要记录日志到`glog`，并返回错误码`pb.StatusCode`

因为db已经是底层服务，所以不再依赖其他rpc。
