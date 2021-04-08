package main

import (
    "encoding/json"
    "github.com/tatoshko/tbot/assets"
    config2 "github.com/tatoshko/tbot/config"
    "github.com/tatoshko/tbot/core"
    "io/ioutil"
    "net/http"
    "os"
)

var err error
var PORT = os.Getenv("PORT")

func main() {
    var jsonFile *os.File

    if jsonFile, err = os.Open("config.json"); err == nil {
        defer jsonFile.Close()

        var config config2.Config
        bytes, _ := ioutil.ReadAll(jsonFile)

        if err = json.Unmarshal(bytes, &config); err != nil {
            panic(err)
        }

        assets.InitBox()

        go http.ListenAndServe("0.0.0.0:" + PORT, nil)
        go core.StartBot(config.Token, config.Hook)
    } else {
        panic(err)
    }

}
