package handlerReview

import (
    "database/sql"
    "embercat/bot/handlerReview/service"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Next(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    var reviewItem *service.FrogReviewItem

    logger := getLogger("NEXT")
    chatID := update.Message.Chat.ID

    frogReviewService := service.NewFrogReviewService(update.Message.From.ID)

    if reviewItem, err = frogReviewService.Next(); err != nil {
        if err == sql.ErrNoRows {
            msg := tgbotapi.NewMessage(chatID, "Ревью закончено, или не начато")
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
            return
        }

        logger(err.Error())
        return
    }

    msg := tgbotapi.NewPhotoShare(chatID, reviewItem.PhotoId)
    msg.Caption = reviewItem.PhotoId

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
