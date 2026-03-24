package handlerTurbo

import (
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "golang.org/x/text/feature/plural"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

func Show(bot *tba.BotAPI, update tba.Update) {
    var err error

    logger := getLogger("Show")

    chatID := update.Message.Chat.ID
    args := update.Message.CommandArguments()

    var liner Liner
    if liner, err = NewLinerFromString(args); err != nil {
        msg := tba.NewMessage(chatID, fmt.Sprintf("Неверный номер вкладыша <b>%s</b>", args))
        msg.ParseMode = tba.ModeHTML

        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    var b []byte
    if b, err = liner.ToPicture(); err != nil {
        logger(err.Error())
        return
    }

    msgp := tba.NewPhoto(update.Message.Chat.ID, tba.FileBytes{Name: liner.ToString(), Bytes: b})
    if _, err := bot.Send(msgp); err != nil {
        logger(err.Error())
    }

    var collection Collection
    if collection, err = LoadCollection(int64(update.Message.From.ID)); err != nil {
        logger(err.Error())

        msg := tba.NewMessage(chatID, fmt.Sprintf("Что-то не так с твоей колекцией/n%s", err.Error()))
        msg.ParseMode = tba.ModeHTML
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    score := collection.ScoreOf(liner)

    if err = message.Set(language.Russian, "В твоей коллекции %d вкладышей",
        plural.Selectf(1, "%d",
            "=0", "У тебя <b>нет</b> вкладыша",
            plural.One, "У тебя пока <b>только один</b> вкладыш",
            plural.Few, "В твоей коллекции <b>%d</b> вкладыша",
            plural.Many, "В твоей коллекции <b>%d</b> вкладышей",
        ),
    ); err != nil {
        logger(err.Error())
    }

    printer := message.NewPrinter(language.Russian)
    txt := fmt.Sprintf("%s <b>%s</b>", printer.Sprintf("В твоей коллекции %d вкладышей", score), liner.ToString())

    msg := tba.NewMessage(chatID, txt)
    msg.ParseMode = tba.ModeHTML

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
