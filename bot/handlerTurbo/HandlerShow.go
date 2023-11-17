package handlerTurbo

import (
    redis2 "embercat/redis"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "golang.org/x/text/feature/plural"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

func HandlerShow(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error

    redis := redis2.GetClient()
    if redis == nil {
        return
    }
    defer redis.Close()

    chatID := update.Message.Chat.ID
    args := update.Message.CommandArguments()

    var liner Liner
    if liner, err = NewLinerFromString(args); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Неверный номер вкладыша <b>%s</b>", args))
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err = bot.Send(msg); err != nil {
            logErr(err)
        }
        return
    }

    var b []byte
    if b, err = liner.GetPicture(); err != nil {
        logErr(err)
        return
    }

    msgp := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: liner.ID, Bytes: b})
    if _, err := bot.Send(msgp); err != nil {
        logErr(err)
    }

    var collection Collection
    if collection, err = LoadCollection(redis, int64(update.Message.From.ID)); err != nil {
        logErr(err)

        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Что-то не так с твоей колекцией/n%s", err.Error()))
        msg.ParseMode = tgbotapi.ModeHTML
        if _, err = bot.Send(msg); err != nil {
            logErr(err)
        }
        return
    }

    score := collection.ScoreOf(liner)

    message.Set(language.Russian, "В твоей коллекции %d вкладышей",
        plural.Selectf(1, "%d",
            "=0", "У тебя <b>нет</b> вкладыша",
            plural.One, "У тебя пока <b>только один</b> вкладыш",
            plural.Few, "В твоей коллекции <b>%d</b> вкладыша",
            plural.Many, "В твоей коллекции <b>%d</b> вкладышей",
        ),
    )

    printer := message.NewPrinter(language.Russian)
    message := fmt.Sprintf("%s <b>%s</b>", printer.Sprintf("В твоей коллекции %d вкладышей", score), liner.ID)

    msg := tgbotapi.NewMessage(chatID, message)
    msg.ParseMode = tgbotapi.ModeHTML

    if _, err := bot.Send(msg); err != nil {
        logErr(err)
    }
}
