package handlerTurbo

import (
	"embercat/assets"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

func HandlerShow(api *tgbotapi.BotAPI, update tgbotapi.Update) {
	parts := strings.SplitN(update.Message.Text, " ", 2)

	box := assets.GetBox()
	if b, err := box.Bytes(fmt.Sprintf(TURBO_FILENAME_KEY, parts[1])); err != nil {
		log.Printf("CallbackCollection Bytes error %s", err.Error())
	} else {
		msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: parts[1], Bytes: b})

		if _, err := api.Send(msg); err != nil {
			log.Printf("Delete message error %s", err.Error())
		}
	}
}
