package handlerReview

import (
    "embercat/bot/handlerReview/service"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Stop(bot *tba.BotAPI, update tba.Update) {
    var err error

    logger := getLogger("STOP")
    chatID := update.Message.Chat.ID

    frogReviewService := service.NewFrogReviewService(update.Message.From.ID)

    if err = frogReviewService.Stop(); err != nil {
        logger(err.Error())
    }

    msg := tba.NewMessage(chatID, "Ревью остановленно")

    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
