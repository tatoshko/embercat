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

    pg := pgsql.GetClient()
    q := `select 1 from turbo where userid = $1 and createdAt = now()`

    r, err := pg.Exec(q)

    logger(fmt.Sprintf("%q | %q", r, err))

    return
    //
    //// Only one gum per day
    //todayer := NewTodayer(redis, int64(userID))
    //var dirty bool
    //
    //if dirty, err = todayer.Dirty(); err != nil {
    //    logErr(err)
    //    return
    //}
    //
    //if dirty {
    //    msg := tgbotapi.NewMessage(chatID, "Можно только одну жвачку в день")
    //    if _, err = bot.Send(msg); err != nil {
    //        log.Printf("Roll bot.Send error %s", err.Error())
    //    }
    //
    //    return
    //}
    //
    //// Ok, load collection
    //var collection Collection
    //if collection, err = LoadCollection(redis, int64(userID)); err != nil {
    //    logErr(err)
    //    msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Что-то не так с твоей коллекцией\n%s", err.Error()))
    //    msg.ParseMode = tgbotapi.ModeHTML
    //
    //    return
    //}
    //
    //// Add new liner to collection
    //liner := GetRandomLiner()
    //if _, err := collection.Add(liner); err != nil {
    //    logErr(err)
    //    msg := tgbotapi.NewMessage(
    //        chatID,
    //        fmt.Sprintf("Что-то не так с выдачей новых вкладышей, а тебе почти достался <b>%s</b>", liner.ID),
    //    )
    //    msg.ParseMode = tgbotapi.ModeHTML
    //
    //    if _, err = bot.Send(msg); err != nil {
    //        log.Printf("Roll bot.Send error %s", err.Error())
    //    }
    //}
    //
    //// Load liner's picture
    //if b, err := liner.GetPicture(); err != nil {
    //    logErr(err)
    //} else {
    //    msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: liner.ID, Bytes: b})
    //
    //    if _, err = bot.Send(msg); err != nil {
    //        logErr(err)
    //    }
    //}

}
