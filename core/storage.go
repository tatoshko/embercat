package core

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "strings"
)

var Storage = make(map[string]string)

func handleSet(update tgbotapi.Update, text string) (tgbotapi.Message, error) {
    parts := strings.SplitAfterN(text, " ", 2)
    key, value := parts[0], parts[1]
    Storage[key] = value

    msg := tgbotapi.NewMessage(
        update.Message.Chat.ID,
        fmt.Sprintf("'%s' has benn set to key '%s'", value, key),
    )
    msg.ReplyToMessageID = update.Message.MessageID

    return Bot.Send(msg)
}

func handleGet(update tgbotapi.Update, key string) (tgbotapi.Message, error) {
    value := Storage[key]

    msg := tgbotapi.NewMessage(
        update.Message.Chat.ID,
        fmt.Sprintf("%s", value),
    )
    msg.ReplyToMessageID = update.Message.MessageID

    return Bot.Send(msg)
}
