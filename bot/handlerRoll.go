package bot

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "time"
)

func handleRoll(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("https://animego.org/anime/random?%d", time.Now().Unix()))
    msg.ReplyToMessageID = update.Message.MessageID
    msg.DisableWebPagePreview = true

    if _, err := bot.Send(msg); err != nil {
        log.Printf("handleRoll error %s\n", err.Error())
    }
}
