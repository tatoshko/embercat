package handlerQuote

import (
    "embercat/bot/handlerQuote/service"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Rnd(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("RND")

    chatID := update.Message.Chat.ID

    quoteService := service.NewService(chatID)

    // load rnd quote
    var quote *service.Quote
    if quote, err = quoteService.FindRND(); err != nil {
        logger(err.Error())
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("что-то пошло не так %s", err.Error()))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    go quoteService.AddStat(quote, service.PlaceRNDMsg)

    msg := tgbotapi.NewMessage(chatID, quote.ToString())
    msg.ReplyToMessageID = update.Message.MessageID

    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
