package main

import (
    "encoding/json"
    "github.com/tatoshko/tbot/assets"
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

        assets.InitBox()

        var config core.Config

        bytes, _ := ioutil.ReadAll(jsonFile)

        if err = json.Unmarshal(bytes, &config); err != nil {
            panic(err)
        }

        go http.ListenAndServe("0.0.0.0:" + PORT, nil)

        core.InitBot(config).Watch()
    } else {
        panic(err)
    }

}
