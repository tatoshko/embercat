package core

import (
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "strings"
)

var err error
var Bot *tba.BotAPI

func InitBot(token, hook string, output chan string) {
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

            output <- fmt.Sprintf("Command '%s' with value '%s'", command, value)

            switch command {
            case "set":
                handleSet(update, value)
            case "get":
                handleGet(update, value)
            }
        }

    }
}
