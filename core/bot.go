package core

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "strings"
)

var err error
var Bot *tba.BotAPI

func InitBot(token, hook string) {
    if Bot, err = tba.NewBotAPI(token); err != nil {
        panic(err)
    }

    Bot.SetWebhook(tba.NewWebhook(hook + "/" + token))
    updates := Bot.ListenForWebhook("/" + Bot.Token)

    for update := range updates {
        //output <- update.Message.Text
        text := update.Message.Text

        if strings.HasPrefix(text, "/") {
            text = strings.Trim(text, "/")
            parts := strings.SplitAfterN(text, " ", 2)
            command, value := parts[0], parts[1]

            switch command {
            case "set":
                Bot.Send(tba.NewMessage(update.Message.Chat.ID, "WTF!"))
                if _, err = handleSet(update, value); err != nil {
                    log.Fatal(err)
                }
            case "get":
                if _, err = handleGet(update, value); err != nil {
                    log.Fatal(err)
                }
            }
        }

    }
}
