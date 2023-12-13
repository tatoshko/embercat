package ai

import (
    "embercat/pgsql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
)

func HandlerTextAnalizer(API *tgbotapi.BotAPI, update tgbotapi.Update) {
    if update.Message.Text == "" {
        return
    }

    r := regexp.MustCompile(`(\w+|[А-я])`)
    parts := r.FindAllString(update.Message.Text, -1)

    pg := pgsql.GetClient()

    if tx, err := pg.Begin(); err != nil {
        log.Printf("TextAnalizer error %s", err.Error())
        return
    } else {
        q := `insert into words (word) values ($1) on conflict (word) do update set count = count + 1`

        for _, word := range parts {
            tx.Exec(q, word)
        }

        if err = tx.Commit(); err != nil {
            log.Printf("Text anilizer Commit error %s", err.Error())
        }
    }
}
