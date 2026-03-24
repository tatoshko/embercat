package handlerQuote

import (
    "embercat/bot/handlerQuote/service"
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Rnd(bot *tba.BotAPI, update tba.Update) {
    var err error
    logger := getLogger("RND")

    chatID := update.Message.Chat.ID

    quoteService := service.NewService(chatID)

    // load rnd quote
    var quote *service.Quote
    if quote, err = quoteService.FindRND(); err != nil {
        logger(err.Error())
        msg := tba.NewMessage(chatID, fmt.Sprintf("что-то пошло не так %s", err.Error()))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    go quoteService.AddStat(quote, service.PlaceRNDMsg)

    msg := tba.NewMessage(chatID, quote.ToString())
    msg.ReplyToMessageID = update.Message.MessageID

    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
