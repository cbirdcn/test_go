# 更新、重载（reload）

实时重载和热重载是两个概念。

实时重新加载是在开发环境中，更改源代码后自动重建并重新启动整个应用程序的过程。一般是使用守护进程实现。直观上类似解释型代码的部署方式。

热重载表示在生产环境中，通过仅替换的方式重新加载部分或全部已编译好的二进制文件。目前，Go 本身还不支持热重载，需要自己实现。要让热加载期间保持服务的连接和可用也是个难题。

对短连接服务，可以在负载均衡侧进行流量切换（通过权重、开关等操纵流量转发）。也就是让部分负载停止服务进行更新，另一部分负载承担全部任务，重新启动已更新的负载后让另一部分负载进行停服更新。

## 程序主体更新

比如游戏服务更新的步骤

- golang服务进程运行时监听USR2信号
- 进程收到USR2信号后, 下载新版本的客户端到本地
- fork子进程(启动新版本服务)
- 将上下文, 句柄等信息交到新的子进程
- 新进程开始监听socket请求
- 等待旧服务连接停止

```go
ch := make(chan os.Signal, 10)
signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
for {
    sign := <-ch
    switch sign
    case syscall.SIGUSR2:
        if err := StartNewPro(); err != nil {
            ......
            break
        }
        execSpec := &syscall.ProcAttr{
            Env: os.Environ(),
            Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
        }
        fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
        ......
        process, _ := os.FindProcess(os.Getppid())
        process.Signal(syscall.SIGHUB)
        ......
    }
```

常规的HTTP服务可以直接用`endless`和`grace`。

## 插件更新

运维、C、C++通常采用plugin插件更新的方式，也就是保持程序主体不变，只更新Linux动态链接库.so文件。这样能让主体程序更小，小更新只发布插件即可。

需要具备gcc环境。

程序主体需要打开文件对象，找到需要的函数，再把`interface{}`类型的函数对象断言成函数类型`func()`，最后执行该函数。

正如`plugin`路径下的例子。

因为使用插件的限制很多，比如插件之间依赖、插件耗时过久、插件中定义全局变量、插件之间类型隔离等问题可能无法或很难处理。所以这种更新一般用于简单的业务逻辑，或配置更新。

## 短连接服务用负载均衡进行流量切换

对于短连接服务，生产环境应该使用负载均衡进行流量切换，不应该使用热加载或平滑重启

go 部署在生成环境应该部署二进制，不使用代码部署，所以 `git pull` 等热加载操作不允许的；

使用热加载或平滑重启可能会遇到新进程启动失败的情况，新进程拉起失败，但是老服务已经 stop 或 kill 在退出了无法处理请求，在启动失败的情况下或服务中断。

在传统环境下可以使用 nginx 切换 upstream 实现手动切换流量，在容器环境下可以使用 svc 自动切换流量，实现好下线行为就可以无缝切换流量。

## 开发环境实时加载

有`air`等第三方组件可用

## 参考

[golang实现热更新的常规方式](https://blog.51cto.com/u_2010293/2781898)

[golang生产环境热加载的可行性](https://learnku.com/go/t/64197#reply214693)

[Live Reloading Your Golang Apps](https://thegodev.com/live-reloading/)

[hotswap热加载](https://github.com/edwingeng/hotswap/blob/main/README.zh-CN.md)