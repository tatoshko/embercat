package core

import (
    "encoding/json"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "io"
    "log"
    "net/http"
)

var CATALOG = "https://2ch.hk/b/catalog.json"

func Init() {
    if res, err := http.Get(CATALOG); err != nil {
        log.Fatalln(err)
    } else {
        defer res.Body.Close()

        if body, err := io.ReadAll(res.Body); err != nil {
            log.Fatalln(err)
        } else {
            var board Board

            if err = json.Unmarshal(body, &board); err != nil {
                log.Fatalln(err)
            } else {
                log.Printf("Subject of first thread %q", board.Threads[0].Subject)
            }
        }

    }
}

func handle2ch(bot *tgbotapi.BotAPI, update tgbotapi.Update, data string) {
    Init()
}
