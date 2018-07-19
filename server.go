package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "encoding/json"
)

type Node struct {
    Name string `json:"nodename"`
    Ip  string `json:"nodeip"`
}

func main() {
    started := time.Now()
    http.HandleFunc("/started", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        var node Node
        json.NewDecoder(r.Body).Decode(&node)
        fmt.Println(node.Name)
        data := (time.Since(started)).String() + node.Name + node.Ip
        w.Write([]byte(data))
    })
    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        duration := time.Since(started)
        if duration.Seconds() > 10 {
            w.WriteHeader(500)
            w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
        } else {
            w.WriteHeader(200)
            w.Write([]byte("ok"))
        }
    })
    log.Fatal(http.ListenAndServe(":9191", nil))
}
