package core

import (
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "strings"
)

var Storage = make(map[string]string)

func handleSet(bot *tba.BotAPI, update tba.Update, text string) {
    parts := strings.SplitAfterN(text, " ", 2)
    key, value := parts[0], parts[1]
    Storage[key] = value

    msg := tba.NewMessage(
        update.Message.Chat.ID,
        fmt.Sprintf("'%s' has been set to key '%s'", value, key),
    )
    msg.ReplyToMessageID = update.Message.MessageID

    bot.Send(msg)
}

func handleGet(bot *tba.BotAPI, update tba.Update, key string) {
    log.Printf("Getted from storage %q, by key %s, result %s", Storage, key, Storage[key])

    value, found := Storage[key]

    msg := tba.NewMessage(
        update.Message.Chat.ID,
        fmt.Sprintf("Value is %v *%s*", found, value),
    )

    msg.ReplyToMessageID = update.Message.MessageID

    bot.Send(msg)
}
