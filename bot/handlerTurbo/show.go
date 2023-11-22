package handlerTurbo

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "golang.org/x/text/feature/plural"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

func Show(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error

    logger := getLogger("Show")

    chatID := update.Message.Chat.ID
    args := update.Message.CommandArguments()

    var liner Liner
    if liner, err = NewLinerFromString(args); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Неверный номер вкладыша <b>%s</b>", args))
        msg.ParseMode = tgbotapi.ModeHTML

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

    msgp := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: liner.ToString(), Bytes: b})
    if _, err := bot.Send(msgp); err != nil {
        logger(err.Error())
    }

    var collection Collection
    if collection, err = LoadCollection(int64(update.Message.From.ID)); err != nil {
        logger(err.Error())

        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Что-то не так с твоей колекцией/n%s", err.Error()))
        msg.ParseMode = tgbotapi.ModeHTML
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

    msg := tgbotapi.NewMessage(chatID, txt)
    msg.ParseMode = tgbotapi.ModeHTML

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
