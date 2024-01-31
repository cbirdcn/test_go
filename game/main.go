package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"

	"test/game/db/dbServer"
	"test/game/link/linkClient"
	"test/game/link/linkServer"
	"test/game/logic/logicServer"
)

func cmdDb(args []string) {
	glog.Info("Starting db server")
	dbServer.Main()
}

func cmdLogic(args []string) {
	glog.Info("Starting logic server")
	logicServer.Main()
}

func cmdLinkServer(args []string) {
	glog.Info("Starting link server")
	linkServer.Main()
}

func cmdLinkClient(args []string) {
	glog.Info("Starting link client")
	linkClient.Main()
}

func usage() {
	fmt.Printf("Usage: %s command args...\n", os.Args[0])
	fmt.Println("Available commands are: db, logic, linkServer, linkClient")
}

// 已废弃
// 原因1：单入口的方式虽然方便运行，但是glog收集的日志都会用main作为文件名，导致日志混乱。
// 原因2：每个单独运行的进程都作为main的子模块了，需要提供一个Main()方法给main模块调用，并不利于解耦开发和部署。
func main() {
	// glog默认的logtostderr是false，日志记录到/tmp/中
	// flag.Set("alsologtostderr", "true")
	// flag.Set("log_dir", "./log")
	flag.Parse()
	defer glog.Flush()

	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	// http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 10

	switch args[0] {
	case "db":
		cmdDb(args[1:])
	case "logic":
		cmdLogic(args[1:])
	case "linkServer":
		cmdLinkServer(args[1:])
	case "linkClient":
		cmdLinkClient(args[1:])
	default:
		usage()
		glog.Fatalf("Unknown command %q", args[0])
	}
}
