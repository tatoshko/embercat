package handlerQuote

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Add(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("ADD")
    chatID := update.Message.Chat.ID
    service := NewService()

    quote := NewQuoteFromMessage(update.Message)

    msg := tgbotapi.NewMessage(chatID, "")

    if err = service.Add(quote); err != nil {
        logger(err.Error())
        msg.Text = fmt.Sprintf("Что-то не так: %s", err.Error())
    } else {
        msg.Text = fmt.Sprintf("Записал цитатку:\n%s", quote.ToString())
    }

    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
