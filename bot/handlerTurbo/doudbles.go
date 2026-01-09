package handlerTurbo

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Double(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("DOUBLES")

    chatID := update.Message.Chat.ID
    userID := update.Message.From.ID

    var collection Collection
    if collection, err = LoadCollection(int64(userID)); err != nil {
        msg := tgbotapi.NewMessage(chatID, "Что-то не так с твоей коллекцией")
        if _, err = bot.Send(msg); err != nil {
            logger("bot send error: %s", err.Error())
        }
        return
    }

    duplicates := collection.GetDuplicates()

    if len(duplicates) < 1 {
        msg := tgbotapi.NewMessage(chatID, "В твоей коллекции нет дубликатов")
        if _, err = bot.Send(msg); err != nil {
            logger("bot send error: %s", err.Error())
        }
    } else {
        var text string

        for _, liner := range duplicates {
            text += fmt.Sprintf("#%-15d(%d) ", liner.ID, liner.Count)
        }

        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("<code>%s</code>", text))
        msg.ParseMode = tgbotapi.ModeHTML
        if _, err = bot.Send(msg); err != nil {
            logger("bot send error: %s", err.Error())
        }
    }
}
