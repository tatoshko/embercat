package handlerTurbo

import (
    "embercat/pgsql"
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Roll(bot *tba.BotAPI, update tba.Update) {
    var err error
    logger := getLogger("ROLL")

    chatID := update.Message.Chat.ID
    userID := update.Message.From.ID

    // Only one gum per day
    pg := pgsql.GetClient()
    q := `select 1 from turbo where userid = $1 and date(createdat) = current_date`
    r, _ := pg.Exec(q, userID)

    if exists, err := r.RowsAffected(); err != nil {
        logger("can't check exists: ", err.Error())
    } else if exists > 0 {
        msg := tba.NewMessage(chatID, "Можно только одну жвачку в день.")
        if _, err := bot.Send(msg); err != nil {
            logger("bot send error", "only one", err.Error())
        }

        return
    }

    // Load collection
    var collection Collection
    if collection, err = LoadCollection(int64(userID)); err != nil {
        logger("collection load error: ", err.Error())
        return
    }

    // Add new liner to collection
    liner := GetRandomLiner()
    if _, err := collection.Add(liner); err != nil {
        logger("add liner error: ", err.Error())
        msg := tba.NewMessage(
            chatID,
            fmt.Sprintf("Что-то не так с выдачей новых вкладышей, а тебе почти достался <b>%s</b>", liner.ToString()),
        )
        msg.ParseMode = tba.ModeHTML

        if _, err = bot.Send(msg); err != nil {
            logger("bot send error", "add liner", err.Error())
        }
    }

    // Load liner's picture
    if b, err := liner.ToPicture(); err != nil {
        logger("liner to pic error: ", err.Error())
    } else {
        msg := tba.NewPhoto(update.Message.Chat.ID, tba.FileBytes{Name: liner.ToString(), Bytes: b})

        if _, err = bot.Send(msg); err != nil {
            logger("bot send error", "photo upload", err.Error())
        }
    }

}
