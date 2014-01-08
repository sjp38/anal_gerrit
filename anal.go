package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
)

const (
    GerritURL = "https://android-review.googlesource.com"
    FirstSortKey = "002a4ecd00011ff0"
)

func pindent(depth int) {
    for i := 0; i < depth; i++ {
        fmt.Print("    ")
    }
}

func printRawMarshalled(infos interface{}, depth int) {
    switch vinfos := infos.(type) {
    case map[string]interface{}:
        pindent(depth)
        fmt.Println("{")
        for k, v := range vinfos {
            switch vv := v.(type) {
            case string:
        pindent(depth)
                fmt.Printf("%v: (string) - %q\n", k, vv)
            case int:
        pindent(depth)
                fmt.Printf("%v: (int) - %q\n", k, vv)
            default:
        pindent(depth)
                fmt.Printf("%v: (not string, neither int)\n", k)
                printRawMarshalled(v, depth + 1)
            }
        }
        pindent(depth)
        fmt.Printf("}\n")
    case []interface{}:
        pindent(depth)
        fmt.Println("[")
        for k, v := range vinfos {
            switch vv := v.(type) {
            case string:
        pindent(depth)
                fmt.Printf("%vth: (string) - %q\n", k, vv)
            case int:
        pindent(depth)
                fmt.Printf("%vth: (int) - %q\n", k, vv)
            default:
        pindent(depth)
                fmt.Printf("%vth: (not string, neither int)\n", k)
                printRawMarshalled(v, depth + 1)
            }
        }
        pindent(depth)
        fmt.Printf("]\n")
    }
}

func fetchChanges(status string, next_sort_key string) interface{} {
    dest_url := fmt.Sprintf("%s/changes/?q=status:%s&n=25&O=1&N=%s",
            GerritURL, status, next_sort_key)
    resp, err := http.Get(dest_url)
    if err != nil {
        log.Fatal("fetch from net fail", err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal("read fail", err)
    }

    index := bytes.IndexByte(body, '\n')
    body = body[index:]

    var infos interface{}
    err = json.Unmarshal(body, &infos)
    if err != nil {
        log.Fatal("unmarshalling fail", err)
    }
    return infos
}

func main() {
    printRawMarshalled(fetchChanges("open", FirstSortKey), 0)
}
