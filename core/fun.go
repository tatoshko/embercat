package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func handleThread(bot *tgbotapi.BotAPI, data string, chatID int64, update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(chatID, "Обсуждение закрыто. Разговаривайте про что-нибудь другое. Спасибо")
    bot.Send(msg)
}
