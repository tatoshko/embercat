package handlerReview

import (
    "database/sql"
    "embercat/bot/handlerReview/service"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "strings"
)

const (
    CBFRRemove = "frrem"
    CBFRStay   = "frstay"
)

func Next(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    next(bot, update.Message.Chat.ID, update.Message.From.ID)
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
    chatID := query.Message.Chat.ID

    itemId := strings.Split(strings.TrimPrefix(query.Data, fmt.Sprintf("/%s ", CBFRRemove)), " ")[0]

    frogReviewService := service.NewFrogReviewService(userID)
    if reviewItem, err := frogReviewService.FindById(itemId); err != nil {
        if err == sql.ErrNoRows {
            msg := tgbotapi.NewMessage(chatID, "Такой жабы уже нет")
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
        } else {
            logger(err.Error())
        }
    } else {
        if err = frogReviewService.Reject(reviewItem); err != nil {
            logger(err.Error())
        } else {
            msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Минус жаба %s", reviewItem.PhotoId))
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }

            del := tgbotapi.NewDeleteMessage(chatID, query.Message.MessageID)
            if _, err := bot.Send(del); err != nil {
                logger(err.Error())
            }

            next(bot, chatID, userID)
        }
    }
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
    chatID := query.Message.Chat.ID
    itemId := strings.Split(strings.TrimPrefix(query.Data, fmt.Sprintf("/%s ", CBFRStay)), " ")[0]

    frogReviewService := service.NewFrogReviewService(userID)
    if reviewItem, err := frogReviewService.FindById(itemId); err != nil {
        if err == sql.ErrNoRows {
            msg := tgbotapi.NewMessage(chatID, "Такой жабы уже нет")
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
        } else {
            logger(err.Error())
        }
    } else {
        if err = frogReviewService.Approve(reviewItem); err != nil {
            logger(err.Error())
        } else {
            msg := tgbotapi.NewMessage(chatID, "Пусть живет")
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }

            del := tgbotapi.NewDeleteMessage(chatID, query.Message.MessageID)
            if _, err := bot.Send(del); err != nil {
                logger(err.Error())
            }

            next(bot, chatID, userID)
        }
    }
}

func next(bot *tgbotapi.BotAPI, chatID int64, userID int) {
    var err error
    var reviewItem *service.FrogReviewItem

    logger := getLogger("NEXT")

    frogReviewService := service.NewFrogReviewService(userID)

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
        tgbotapi.NewInlineKeyboardButtonData("☦ Remove", fmt.Sprintf("/%s %s", CBFRRemove, reviewItem.Id)),
        tgbotapi.NewInlineKeyboardButtonData("✅ Stay", fmt.Sprintf("/%s %s", CBFRStay, reviewItem.Id)),
    ))

    msg := tgbotapi.NewPhotoShare(chatID, reviewItem.PhotoId)
    msg.Caption = reviewItem.PhotoId
    msg.ReplyMarkup = keyboard

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
