package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func handlerYesNo(api *tba.BotAPI, update tba.Update)  {
    msg := tba.NewMessage(update.Message.Chat.ID, "Yes/No")

    keyboard := tba.InlineKeyboardMarkup{
        InlineKeyboard: [][]tba.InlineKeyboardButton{
            {
                tba.NewInlineKeyboardButtonData("Yes", "y"),
                tba.NewInlineKeyboardButtonData("No", "n"),
            },
        },
    }

    msg.ReplyMarkup = keyboard

    if _, err := api.Send(msg); err != nil {
        log.Printf("handlerYesNo error %s\n", err.Error())
    } else {
        registerCallback(update.Message.MessageID, func(api *tba.BotAPI, update tba.Update) {

        })
    }
}
