# logic逻辑服务

## 简介

作为rpc server向外提供服务，调用方可能不止一个，比如gmt（http）、wsServer（ws）。

另外因为wsServer只有转发功能，所以logic需要兼具响应包装功能。

所以，logic需要根据"调用方-路由"做出不同的逻辑处理。

因为要做成高并发以及大量玩家同服，所以用mongodb集群作为存储。这就需要在logic启动前保证db service启动。

目前只实现了基础功能

## 说明
