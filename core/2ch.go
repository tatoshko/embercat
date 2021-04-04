package core

import (
    "encoding/json"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "io"
    "log"
    "net/http"
    "regexp"
)

var CATALOG = "https://2ch.hk/b/catalog.json"
var db []File = make([]File, 1000)

func parseThreads() []Thread {
    var threads []Thread = make([]Thread, 10)
    var matcher = regexp.MustCompile(`(webm|вебм|цуин|mp4)`)

    var board Board

    if err = lookup(CATALOG, &board); err != nil {
        log.Fatalln(err)
    } else {
        for _, thread := range board.Threads {
            if matcher.MatchString(thread.Subject) || matcher.MatchString(thread.Comment){
                threads = append(threads, thread)
            }
        }
    }

    return threads
}

func handle2ch(bot *tgbotapi.BotAPI, update tgbotapi.Update, data string) {
    for _, thread := range parseThreads() {

        if err = lookup(makeThreadPath(thread.Num), &thread); err != nil {
            log.Fatalln(err)
        } else {
            for _, post := range thread.Posts {
                for _, file := range post.Files {
                    if file.Type == MP4 {
                        db = append(db, file)
                    }
                }
            }
        }
    }

    msg := tgbotapi.NewVideoUpload(update.Message.Chat.ID, db[0].Path)
    bot.Send(msg)
}

func lookup(path string, data interface{}) (err error) {
    if res, err := http.Get(path); err != nil {
        return
    } else {
        defer res.Body.Close()

        if body, err := io.ReadAll(res.Body); err != nil {
            return
        } else {
            if err = json.Unmarshal(body, &data); err != nil {
                return
            }
        }
    }
}

func makeThreadPath(num string) string {
    return fmt.Sprintf("https://2ch.hk/b/res/%s.html", num)
}