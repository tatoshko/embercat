package core

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
)

var err error
var bot *tba.BotAPI

func InitBot(token, hook string, output chan string) {
    if bot, err = tba.NewBotAPI(token); err != nil {
        panic(err)
    }

    bot.SetWebhook(tba.NewWebhook(hook + "/" + token))
    updates := bot.ListenForWebhook("/" + bot.Token)

    for update := range updates {
        output <- update.Message.Text
    }
}
