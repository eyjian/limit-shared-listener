// Writed by yijian on 2021/01/14
package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

var (
    help = flag.Bool("h", false, "Display a help message and exit.")
    port1 = flag.Uint("port1", 20210, "The first port of server.")
    port2 = flag.Uint("port2", 20211, "The second port of server.")
)

func main()  {
    flag.Parse()
    if *help {
        flag.Usage()
        os.Exit(1)
    }

    httpClient1 := http.Client{ Timeout: time.Duration(1)*time.Second}
    httpClient2 := http.Client{ Timeout: time.Duration(1)*time.Second}
    // 不加上“http:”会报错误："first path segment in URL cannot contain colon"
    resp1, err1 := httpClient1.Get(fmt.Sprintf("http://127.0.0.1:%d/hello?hello1", *port1))
    if err1 != nil {
        fmt.Printf("First get error: %s\n", err1.Error())
    } else {
        fmt.Printf("First get ok\n")
        body, err := ioutil.ReadAll(resp1.Body)
        if err != nil {
            fmt.Printf("First read error: %s\n", err.Error())
        } else {
            fmt.Println(string(body))
        }
    }
    resp2, err2 := httpClient2.Get(fmt.Sprintf("http://127.0.0.1:%d/hello?hello2", *port2))
    if err2 != nil {
        fmt.Printf("Second get error: %s\n", err2.Error())
    }else {
        fmt.Printf("Second get ok\n")
        body, err := ioutil.ReadAll(resp2.Body)
        if err != nil {
            fmt.Printf("Second read error: %s\n", err.Error())
        } else {
            fmt.Println(string(body))
        }
    }
    time.Sleep(time.Duration(10) * time.Second)
}
