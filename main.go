package main

import (
    "embercat/assets"
    "embercat/bot"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

var err error
var HOST = ""
var PORT = ":3001"

type Config struct {
    Name  string `json:"name"`
    Token string `json:"token"`
    Hook  string `json:"hook"`
}

func main() {
    var jsonFile *os.File

    if jsonFile, err = os.Open("config.json"); err == nil {
        defer jsonFile.Close()

        var config Config
        bytes, _ := io.ReadAll(jsonFile)

        if err = json.Unmarshal(bytes, &config); err != nil {
            panic(err)
        }

        assets.InitBox()

        go (func() {
            addr := fmt.Sprintf("%s%s", HOST, PORT)

            log.Printf("Bind to: [%s]", addr)

            if err := http.ListenAndServe(addr, nil); err != nil {
                log.Fatalln(err.Error())
            }
        })()

        bot.Start(config.Name, config.Token, config.Hook)
    } else {
        panic(err)
    }

}
