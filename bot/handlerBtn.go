package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

const SOME_DATA = "some_data"

func handlerBtn(api *tba.BotAPI, update tba.Update)  {
    msg := tba.NewMessage(update.Message.Chat.ID, "Это кнопка для нажимания на неё")

    keyboard := tba.InlineKeyboardMarkup{
        InlineKeyboard: [][]tba.InlineKeyboardButton{
            {
                tba.NewInlineKeyboardButtonData("Жми уже", SOME_DATA),
            },
        },
    }

    msg.ReplyMarkup = keyboard

    if _, err := api.Send(msg); err != nil {
        log.Printf("handlerBtn error %s\n", err.Error())
    } else {
        delMsg := tba.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
        if _, err := api.Send(delMsg); err != nil {
            log.Printf("Delete message error %s\n", err.Error())
        }

        registerCallback(SOME_DATA, func(api *tba.BotAPI, update tba.Update) {
            callback := tba.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
            if _, err := api.AnswerCallbackQuery(callback); err != nil {
                log.Printf("Callback error %s", err.Error())
            } else {
                delMsg := tba.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

                if _, err := api.Send(delMsg); err != nil {
                    log.Printf("Delete message error %s", err.Error())
                }

                msg := tba.NewMessage(update.CallbackQuery.Message.Chat.ID, "А если бы рвануло? Не жми на все кнопки подряд")

                if _, err := api.Send(msg); err != nil {
                    log.Printf("Delete message error %s", err.Error())
                }

                unregisterCallback(SOME_DATA)
            }
        })
    }
}
