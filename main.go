package main

import (
    "embercat/assets"
    "embercat/bot"
    "embercat/pgsql"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

var err error

type ServerConfig struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

type Config struct {
    Bot    bot.Config   `json:"bot"`
    Pg     pgsql.Config `json:"pg"`
    Server ServerConfig `json:"server"`
}

func main() {
    config := ReadConfig()

    assets.InitBox()
    go initHttpServer(config.Server)
    go pgsql.Init(config.Pg)
    bot.Start(config.Bot)

}

func ReadConfig() (config Config) {
    var jsonFile *os.File
    if jsonFile, err = os.Open("config.json"); err != nil {
        panic(err)
    }
    defer jsonFile.Close()

    bytes, _ := io.ReadAll(jsonFile)

    if err = json.Unmarshal(bytes, &config); err != nil {
        panic(err)
    }

    return
}

func initHttpServer(config ServerConfig) {
    addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

    log.Printf("Bind to: [%s]", addr)

    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalln(err.Error())
    }
}
