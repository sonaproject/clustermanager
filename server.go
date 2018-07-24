package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "strconv"
    "time"
)

var path = "/Users/novalagung/Documents/temp/test.txt"
var etag = 0;
type Node struct {
    Name string `json:"nodename"`
    Ip  string `json:"nodeip"`
}
var modTime time.Time
var now time.Time

func main() {
    fmt.Println("Build 0.11")
    var ready bool
    var nodeStore map[string]string
    nodeStore = make(map[string]string)
    ready = false
    now = time.Now()
    http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        var node Node
        json.NewDecoder(r.Body).Decode(&node)
        fmt.Println(node.Name + " has joind with IP " + node.Ip)
        data := node.Name + ":" + node.Ip
        w.Write([]byte(data))
        nodeStore[node.Name] = node.Ip
        if len(nodeStore) > 2 {
          ready = true
        }
        etag++
        modTime = now
        now = time.Now()
    })
    http.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
        var data string
        w.WriteHeader(200)
        data = ""
        for key, value := range nodeStore {
          data += key + ":" + value + "\n"
        }
        w.Write([]byte(data))
    })
    http.HandleFunc("/leave", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        var node Node
        json.NewDecoder(r.Body).Decode(&node)
        fmt.Println(node.Name + " has left")
        data := node.Name + " has left"
        w.Write([]byte(data))
        delete( nodeStore, node.Name)
        etag++
        modTime = now
        now = time.Now()
    })
    http.HandleFunc("/cluster.json", func(w http.ResponseWriter, r *http.Request) {
        var data string
        var header string
        var footer string
        var nodeInfo string
        var memberInfo string
        var middle string
        var etagHeader string
        data = ""
        if (ready && len(nodeStore) > 1) || (!ready && len(nodeStore) > 2) {
          header = "{\"nodes\": ["
          middle = "],\"name\": 12345678, \"partitions\": [{\"id\": 1, \"members\": ["
          footer = "]}]}"
          nodeInfo = ""
          var numMember = len (nodeStore)
          var i = 0;
          for key, value := range nodeStore {
            i++
            nodeInfo += "{ \"ip\" : \"" + value + "\", \"id\": \"" + key + "\", \"port\" : 9876 }\n"
            memberInfo += "\"" + key + "\""
            if i < numMember {
              nodeInfo += ","
              memberInfo += ","
            }
          }
          data = header + nodeInfo + middle + memberInfo + footer
        }
        etagHeader =  "\"" + strconv.Itoa(etag) + "\""
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Connection", "close")
        w.Header().Set("ETag", etagHeader)
        w.Header().Set("Accept-Ranges", "bytes")
        w.Header().Set("Last-Modified", modTime.UTC().Format(http.TimeFormat))
        w.Write([]byte(data))
    })
    log.Fatal(http.ListenAndServe(":9191", nil))
}
