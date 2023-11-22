package handlerTurbo

import (
    "embercat/pgsql"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Roll(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("ROLL")

    chatID := update.Message.Chat.ID
    userID := update.Message.From.ID

    // Only one gum per day
    pg := pgsql.GetClient()
    q := `select 1 from turbo where userid = $1 and date(createdAt) = current_date`
    r, _ := pg.Exec(q, userID)

    if exists, err := r.RowsAffected(); err != nil {
        logger(err.Error())
    } else if exists > 0 {
        msg := tgbotapi.NewMessage(chatID, "Можно только одну жвачку в день.")
        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    }

    // Load collection
    var collection Collection
    if collection, err = LoadCollection(int64(userID)); err != nil {
        logger(err.Error())
        return
    }

    // Add new liner to collection
    liner := GetRandomLiner()
    if _, err := collection.Add(liner); err != nil {
        logger(err.Error())
        msg := tgbotapi.NewMessage(
            chatID,
            fmt.Sprintf("Что-то не так с выдачей новых вкладышей, а тебе почти достался <b>%s</b>", liner.ID),
        )
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }

    // Load liner's picture
    if b, err := liner.ToPicture(); err != nil {
        logger(err.Error())
    } else {
        msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: liner.ToString(), Bytes: b})

        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }

}
