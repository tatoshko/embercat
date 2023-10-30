package bot

import (
	"embercat/assets"
	"embercat/bot/callbacks"
	"fmt"
	tba "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func handlerBtn(api *tba.BotAPI, update tba.Update) {
	MessageID := fmt.Sprint("%d", update.Message.MessageID)

	delMsg := tba.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)

	if _, err := api.Send(delMsg); err != nil {
		log.Printf("Delete message error %s", err.Error())
	}

	box := assets.GetBox()
	b, err := box.Bytes("prodolbat.png")
	if err != nil {
		log.Println(err)
		return
	}

	msg := tba.NewPhotoUpload(update.Message.Chat.ID, tba.FileBytes{Name: "prodolbatel9000", Bytes: b})

	keyboard := tba.InlineKeyboardMarkup{
		InlineKeyboard: [][]tba.InlineKeyboardButton{
			{
				tba.NewInlineKeyboardButtonData("Жми уже", MessageID),
			},
		},
	}

	msg.ReplyMarkup = keyboard

	if _, err := api.Send(msg); err != nil {
		log.Printf("handlerBtn error %s\n", err.Error())
	} else {
		callbacks.RegisterCallback(MessageID, func(api *tba.BotAPI, update tba.Update) {
			callback := tba.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := api.AnswerCallbackQuery(callback); err != nil {
				log.Printf("Callback error %s", err.Error())
			} else {
				delMsg := tba.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

				if _, err := api.Send(delMsg); err != nil {
					log.Printf("Delete message error %s", err.Error())
				}

				msg := tba.NewMessage(update.CallbackQuery.Message.Chat.ID, "Успешно продолбано")

				if _, err := api.Send(msg); err != nil {
					log.Printf("Delete message error %s", err.Error())
				}

				callbacks.UnregisterCallback(MessageID)
			}
		})
	}
}
