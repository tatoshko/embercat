package bot

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func handleThread(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*Обсуждение закрыто*. Разговаривайте про что-нибудь другое. Спасибо")
    msg.ParseMode = tgbotapi.ModeMarkdown

    if _, err := bot.Send(msg); err != nil {
        log.Printf("handleThread error %s\n", err.Error())
    }
}
