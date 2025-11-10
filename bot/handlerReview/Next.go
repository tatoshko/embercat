package handlerReview

import (
    "database/sql"
    "embercat/bot/handlerReview/service"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "strings"
)

const (
    CBFRRemove = "frremove"
    CBFRStay   = "frstay"
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

    keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("☦ Remove", fmt.Sprintf("/%s %s", CBFRRemove, reviewItem.FrogId)),
        tgbotapi.NewInlineKeyboardButtonData("✅ Stay", fmt.Sprintf("/%s %s", CBFRStay, reviewItem.FrogId)),
    ))

    msg := tgbotapi.NewPhotoShare(chatID, reviewItem.PhotoId)
    msg.Caption = reviewItem.PhotoId
    msg.ReplyMarkup = keyboard

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}

func CallbackRemove(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger := getLogger("CALLBACK REMOVE")

    query := update.CallbackQuery
    callback := tgbotapi.NewCallback(query.ID, query.Data)

    if _, err := bot.AnswerCallbackQuery(callback); err != nil {
        logger(err.Error())
        return
    }

    userID := query.From.ID
    data := strings.Split(strings.TrimLeft(query.Data, fmt.Sprintf("/%s ", CBFRRemove)), " ")

    logger(fmt.Sprintf("%v %v", userID, data))
}

func CallbackStay(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger := getLogger("CALLBACK STAY")

    query := update.CallbackQuery
    callback := tgbotapi.NewCallback(query.ID, query.Data)

    if _, err := bot.AnswerCallbackQuery(callback); err != nil {
        logger(err.Error())
        return
    }

    userID := query.From.ID
    data := strings.Split(strings.TrimLeft(query.Data, fmt.Sprintf("/%s ", CBFRStay)), " ")

    logger(fmt.Sprintf("%v %v", userID, data))
}
