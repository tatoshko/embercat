package main

import (
    "encoding/json"
    "fmt"
    "github.com/tatoshko/tbot/core"
    "io/ioutil"
    "net/http"
    "os"
)

var err error
var output = make(chan string)
var PORT = os.Getenv("PORT")

func main() {
    var jsonFile *os.File

    if jsonFile, err = os.Open("config.json"); err == nil {
        defer jsonFile.Close()

        var config core.Config

        bytes, _ := ioutil.ReadAll(jsonFile)

        if err = json.Unmarshal(bytes, &config); err != nil {
            panic(err)
        }

        go core.InitBot(config.Token, config.Hook, output)

        if err = http.ListenAndServe("0.0.0.0:" + PORT, nil); err != nil {
            panic(err)
        }

        select {
        case msg := <- output:
            fmt.Println(msg)
        }
    } else {
        panic(err)
    }

}