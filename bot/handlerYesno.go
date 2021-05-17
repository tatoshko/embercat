package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

const SOME_DATA = "some_data"

func handlerYesNo(api *tba.BotAPI, update tba.Update)  {
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
        log.Printf("handlerYesNo error %s\n", err.Error())
    } else {
        registerCallback(SOME_DATA, func(api *tba.BotAPI, update tba.Update) {
            msg := tba.NewMessage(update.CallbackQuery.Message.Chat.ID, "А если бы рвануло? Не жми на все кнопки подряд")
            api.Send(msg)
        })
    }
}
