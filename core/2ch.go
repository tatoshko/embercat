package core

import (
    "encoding/json"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "io"
    "log"
    "net/http"
    "net/url"
    "regexp"
)

var CATALOG = "https://2ch.hk/b/catalog.json"

func parseThreads() []Thread {
    var threads []Thread = make([]Thread, 10)
    var matcher = regexp.MustCompile(`(webm|вебм|цуин|mp4)`)

    var board Board

    if err = lookup(CATALOG, &board); err != nil {
        log.Fatalln(err)
    } else {
        for _, thread := range board.Threads {
            if matcher.MatchString(thread.Subject) || matcher.MatchString(thread.Comment) {
                threads = append(threads, thread)
            }
        }
    }

    return threads
}

func handle2ch(bot *tgbotapi.BotAPI, update tgbotapi.Update, data string) {
    for _, thread := range parseThreads() {
        if thread.Num == "" {
            continue
        }

        threadPath := makeThreadPath(thread.Num)

        var board Board

        if err = lookup(threadPath, &board); err != nil {
            log.Println(err)
        } else {
            thread := board.Threads[0]
            var db = make([]File, 1)

            for _, post := range thread.Posts {
                for _, file := range post.Files {
                    if file.Type == MP4 && file.Path != "" {
                        db = append(db, file)
                    }
                }
            }

            //for _, file := range db {
            //    log.Printf("File path: %s", file.Path)
            //}

            fileURL := url.URL{
                Scheme:      "https",
                Host:        "2ch.hk",
                Path:        "/b/src/243772897/16175338978370.mp4",
            }

            if res, err := http.Get(fileURL.String()); err == nil {
                defer res.Body.Close()

                if body, err := io.ReadAll(res.Body); err == nil {
                    msg := tgbotapi.NewVideoUpload(update.Message.Chat.ID, body)
                    bot.Send(msg)
                } else {
                    log.Fatalln(err)
                }
            } else {
                log.Fatalln(err)
            }
        }
    }
}

func lookup(path string, data interface{}) error {
    if res, err := http.Get(path); err != nil {
        return err
    } else {
        defer res.Body.Close()

        if body, err := io.ReadAll(res.Body); err != nil {
            return err
        } else {
            if err = json.Unmarshal(body, &data); err != nil {
                return err
            }
        }
    }

    return nil
}

func makeThreadPath(num string) string {
    return fmt.Sprintf("https://2ch.hk/b/res/%s.json", num)
}
