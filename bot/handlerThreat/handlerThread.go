package handlerThreat

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
)

func HandleThread(bot *tba.BotAPI, update tba.Update) {
    msg := tba.NewMessage(update.Message.Chat.ID, "*Обсуждение закрыто*. Разговаривайте про что-нибудь другое. Спасибо")
    msg.ParseMode = tba.ModeMarkdown

    if _, err := bot.Send(msg); err != nil {
        log.Printf("handleThread error %s\n", err.Error())
    }
}
