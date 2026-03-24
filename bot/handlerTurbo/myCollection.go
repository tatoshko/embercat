package handlerTurbo

import (
    "bytes"
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "golang.org/x/text/feature/plural"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
    "image/jpeg"
    "log"
)

func MyCollection(bot *tba.BotAPI, update tba.Update) {
    var err error

    logger := getLogger("MY COLLECTION")

    chatID := update.Message.Chat.ID
    userID := update.Message.From.ID

    // Load collection and notify user
    var collection Collection
    if collection, err = LoadCollection(int64(userID)); err != nil {
        logger(err.Error())
        msg := tba.NewMessage(chatID, fmt.Sprintf("Что-то не так с твоей коллекцией\n%s", err.Error()))
        msg.ParseMode = tba.ModeHTML

        return
    }

    if err = message.Set(language.Russian, "В твоей коллекции %d вкладышей",
        plural.Selectf(1, "%d",
            "=0", "У тебя <b>нет</b> вкладышей",
            "=1", "У тебя пока <b>только один</b> вкладыш",
            "=2", "В твоей коллекции <b>Два</b> вкладыша",
            "=3", "В твоей коллекции <b>Три</b> вкладыша",
            "=4", "В твоей коллекции <b>Четыре</b> вкладыша",
            "=5", "В твоей коллекции <b>Пять</b> вкладышей",
            plural.One, "В твоей коллекции <b>%d</b> вкладыш",
            plural.Few, "В твоей коллекции <b>%d</b> вкладыша",
            plural.Many, "В твоей коллекции <b>%d</b> вкладышей",
        ),
    ); err != nil {
        logger(err.Error())
    }

    printer := message.NewPrinter(language.Russian)
    result := printer.Sprintf("В твоей коллекции %d вкладышей", collection.Count())
    msg := tba.NewMessage(chatID, result)
    msg.ParseMode = tba.ModeHTML
    if _, err := bot.Send(msg); err != nil {
        log.Printf("HandlerCollection send error %s", err.Error())
    }

    // Generate picture and send to chat
    collectionCanvas := collection.GenerateCollectionPicture()

    buf := new(bytes.Buffer)
    if err := jpeg.Encode(buf, collectionCanvas, nil); err != nil {
        log.Printf("HandlerCollection jpeg.Encode error %s", err.Error())
    }

    msgp := tba.NewDocument(chatID, tba.FileBytes{Bytes: buf.Bytes(), Name: "Collection.jpg"})
    if _, err := bot.Send(msgp); err != nil {
        log.Printf("HandlerCollection bot.Send error %s", err.Error())
    }
}
