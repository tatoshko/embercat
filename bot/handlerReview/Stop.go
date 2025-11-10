package handlerReview

import (
    "embercat/bot/handlerReview/service"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Stop(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error

    logger := getLogger("STOP")
    chatID := update.Message.Chat.ID

    frogReviewService := service.NewFrogReviewService(update.Message.From.ID)

    if err = frogReviewService.Stop(); err != nil {
        logger(err.Error())
    }

    msg := tgbotapi.NewMessage(chatID, "Ревью остановленно")

    if _, err = bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
