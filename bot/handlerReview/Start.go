package handlerReview

import (
    "embercat/bot/handlerReview/service"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error

    logger := getLogger("START")
    frogReviewService := service.NewFrogReviewService(update.Message.From.ID)

    if err = frogReviewService.Start(); err != nil {
        logger(err.Error())
        return
    }

    go Next(bot, update)
}
