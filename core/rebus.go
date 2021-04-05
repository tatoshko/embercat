package core

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func handleRebus(tbot *TBot) CommandHandler {
    return func (bot *tgbotapi.BotAPI, update tgbotapi.Update, value string)  {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ddd")
        if m, err := bot.Send(msg); err == nil {
            id := fmt.Sprintf("id_%d", m.MessageID)

            tbot.RegisterReplay(id, "3d")
        } else {
            log.Fatalln(err)
        }
    }
}
