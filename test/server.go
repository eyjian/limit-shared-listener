// Writed by yijian on 2021/01/14
package main

import (
    "flag"
    "fmt"
    "net"
    "net/http"
    "os"
)
import (
    "github.com/eyjian/limit-shared-listener"
)

var (
    help = flag.Bool("h", false, "Display a help message and exit.")
    port1 = flag.Uint("port1", 20210, "The first port of server.")
    port2 = flag.Uint("port2", 20211, "The second port of server.")
    numConnections = flag.Uint("n", 1, "Maximum number of connections.")
)

func main()  {
    flag.Parse()
    if *help {
        flag.Usage()
        os.Exit(1)
    }

    l1, err1 := net.Listen("tcp", fmt.Sprintf(":%d", *port1))
    l2, err2 := net.Listen("tcp", fmt.Sprintf(":%d", *port2))
    if err1 != nil {
        fmt.Println(err1)
        os.Exit(1)
    }
    if err2 != nil {
        fmt.Println(err2)
        os.Exit(1)
    }

    ll := lsl.NewListenLimiter(int(*numConnections))
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", helloHandler)
    httpServer1 := &http.Server{
        Handler: mux,
    }
    httpServer2 := &http.Server{
        Handler: mux,
    }
    go httpServer1.Serve(lsl.LimitSharedListener(l1, ll))
    httpServer2.Serve(lsl.LimitSharedListener(l2, ll))
    httpServer1.Close()
    httpServer2.Close()
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
    fmt.Println(req.URL.RawQuery)
    w.Write([]byte(req.URL.RawQuery))
}
