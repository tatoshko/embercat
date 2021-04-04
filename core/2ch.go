package core

import (
    "encoding/json"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "io"
    "log"
    "net/http"
    "regexp"
)

var CATALOG = "https://2ch.hk/b/catalog.json"

func parseThreads() []Thread {
    var threads []Thread = make([]Thread, 10)
    var matcher = regexp.MustCompile(`(webm|вебм|цуин|mp4)`)

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
                for _, thread := range board.Threads {
                    if matcher.MatchString(thread.Subject) {
                        threads = append(threads, thread)
                    }
                }
            }
        }
    }

    return threads
}

func handle2ch(bot *tgbotapi.BotAPI, update tgbotapi.Update, data string) {
    for _, thread := range parseThreads() {
        log.Printf("Thread: %s", thread.Subject)
    }
}
