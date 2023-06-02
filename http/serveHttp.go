package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux() // HTTP 请求路由器（多路复用器，Multiplexor）

    rh := http.RedirectHandler("http://www.baidu.com", 307)
    mux.Handle("/foo", rh)

    log.Println("Listening...")
    http.ListenAndServe(":3000", mux)
    // 访问http://localhost:3000/foo
}