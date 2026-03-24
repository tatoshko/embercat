package handlerReview

import (
    "database/sql"
    "embercat/bot/handlerReview/service"
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "strings"
)

const (
    CBFRRemove = "frrem"
    CBFRStay   = "frstay"
)

func Next(bot *tba.BotAPI, update tba.Update) {
    next(bot, update.Message.Chat.ID, update.Message.From.ID)
}

func CallbackRemove(bot *tba.BotAPI, update tba.Update) {
    logger := getLogger("CALLBACK REMOVE")

    query := update.CallbackQuery
    callback := tba.NewCallback(query.ID, query.Data)

    if _, err := bot.Request(callback); err != nil {
        logger(err.Error())
        return
    }

    userID := query.From.ID
    chatID := query.Message.Chat.ID

    itemId := strings.Split(strings.TrimPrefix(query.Data, fmt.Sprintf("/%s ", CBFRRemove)), " ")[0]

    frogReviewService := service.NewFrogReviewService(userID)
    if reviewItem, err := frogReviewService.FindById(itemId); err != nil {
        if err == sql.ErrNoRows {
            msg := tba.NewMessage(chatID, "Такой жабы уже нет")
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
            msg := tba.NewMessage(chatID, fmt.Sprintf("Минус жаба %s", reviewItem.PhotoId))
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }

            del := tba.NewDeleteMessage(chatID, query.Message.MessageID)
            if _, err := bot.Send(del); err != nil {
                logger(err.Error())
            }

            next(bot, chatID, userID)
        }
    }
}

func CallbackStay(bot *tba.BotAPI, update tba.Update) {
    logger := getLogger("CALLBACK STAY")

    query := update.CallbackQuery
    callback := tba.NewCallback(query.ID, query.Data)

    if _, err := bot.Request(callback); err != nil {
        logger(err.Error())
        return
    }

    userID := query.From.ID
    chatID := query.Message.Chat.ID
    itemId := strings.Split(strings.TrimPrefix(query.Data, fmt.Sprintf("/%s ", CBFRStay)), " ")[0]

    frogReviewService := service.NewFrogReviewService(userID)
    if reviewItem, err := frogReviewService.FindById(itemId); err != nil {
        if err == sql.ErrNoRows {
            msg := tba.NewMessage(chatID, "Такой жабы уже нет")
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
            msg := tba.NewMessage(chatID, "Пусть живет")
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }

            del := tba.NewDeleteMessage(chatID, query.Message.MessageID)
            if _, err := bot.Send(del); err != nil {
                logger(err.Error())
            }

            next(bot, chatID, userID)
        }
    }
}

func next(bot *tba.BotAPI, chatID int64, userID int64) {
    var err error
    var reviewItem *service.FrogReviewItem

    logger := getLogger("NEXT")

    frogReviewService := service.NewFrogReviewService(userID)

    if reviewItem, err = frogReviewService.Next(); err != nil {
        if err == sql.ErrNoRows {
            msg := tba.NewMessage(chatID, "Ревью закончено, или не начато")
            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
            return
        }

        logger(err.Error())
        return
    }

    keyboard := tba.NewInlineKeyboardMarkup(tba.NewInlineKeyboardRow(
        tba.NewInlineKeyboardButtonData("☦ Remove", fmt.Sprintf("/%s %s", CBFRRemove, reviewItem.Id)),
        tba.NewInlineKeyboardButtonData("✅ Stay", fmt.Sprintf("/%s %s", CBFRStay, reviewItem.Id)),
    ))

    file := tba.FileID(reviewItem.PhotoId)

    msg := tba.NewPhoto(chatID, file)
    msg.Caption = reviewItem.PhotoId
    msg.ReplyMarkup = keyboard

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
