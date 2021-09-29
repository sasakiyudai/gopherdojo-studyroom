package main

import (
    "fmt"
    "flag"
    "time"
    "net/http"

    "github.com/sasakiyudai/gopherdojo-studyroom/kadai4/sasakiyudai/omikuji"
)

var port int

func init() {
    flag.IntVar(&port, "p", 8080, "port number")
}

func main() {
    flag.Parse()
    o := omikuji.New(time.Now())
    http.HandleFunc("/", o.Handler)
    fmt.Println("おみくじAPIサーバー started...")
    fmt.Printf("type `curl localhost:%d`\n", port)
    http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}