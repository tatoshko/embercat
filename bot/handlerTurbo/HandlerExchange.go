package handlerTurbo

import (
	redis2 "embercat/redis"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func HandlerExchange(api *tgbotapi.BotAPI, update tgbotapi.Update) {
	redis := redis2.GetClient()
	if redis == nil {
		return
	}
	defer redis.Close()

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	localCollectionKey := makeLocalKey(REDIS_KEY_TURBO_COLLECTION, userID)

	if stats, err := redis.ZRangeWithScores(localCollectionKey, 0, -1).Result(); err != nil {
		log.Printf("HandlerCollection ZRangeWithScores error %s", err.Error())
	} else {
		msg := tgbotapi.NewMessage(chatID, "Обменник")

		if update.Message.CommandArguments() == "close" {
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		} else {
			var rows [][]tgbotapi.InlineKeyboardButton
			buttons := make([]tgbotapi.InlineKeyboardButton, 0)
			for i, v := range stats {
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", v.Member), fmt.Sprintf("/ex %s", v.Member)))

				if i%4 == 3 {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons...))
					buttons = make([]tgbotapi.InlineKeyboardButton, 0)
				}
			}

			if len(buttons) > 0 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons...))
			}

			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
		}

		if _, err := api.Send(msg); err != nil {
			log.Printf("HandlerCollection send error %s", err.Error())
		}
	}

}

func CallbackExchange(api *tgbotapi.BotAPI, update tgbotapi.Update) {
	query := update.CallbackQuery

	callback := tgbotapi.NewCallback(query.ID, query.Data)

	if _, err := api.AnswerCallbackQuery(callback); err != nil {
		log.Printf("CallbackExchange answer error %s", err.Error())
	} else {
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, fmt.Sprintf("GOT %s", query.Data))
		if _, err := api.Send(msg); err != nil {
			log.Printf("CallbackExchange send error %s", err.Error())
		}
	}
}
