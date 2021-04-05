package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func handleThread(bot *tgbotapi.BotAPI, update tgbotapi.Update, value string) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Обсуждение закрыто. Разговаривайте про что-нибудь другое. Спасибо")
    bot.Send(msg)
}
