package handlerQuote

import (
    "bytes"
    "embercat/bot/handlerQuote/drawer"
    "embercat/bot/handlerQuote/loader"
    "embercat/bot/handlerQuote/service"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "image/jpeg"
    "log"
)

func Make(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("RND")
    chatID := update.Message.Chat.ID

    // get replay message
    replay := update.Message.ReplyToMessage
    if replay == nil || replay.Photo == nil {
        msg := tgbotapi.NewMessage(chatID, "Ты че пёс, сообщение должно быть с картинкой")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    photos := *replay.Photo

    // get last photoID from replay
    photoID := photos[len(photos)-1].FileID

    // get direct lint to file
    var fileURL string
    if fileURL, err = bot.GetFileDirectURL(photoID); err != nil {
        msg := tgbotapi.NewMessage(chatID, "Не смог получить картинку")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    img, err := loader.LoadPicByURL(fileURL)

    quoteService := service.NewService()

    // load rnd quote
    var quote *service.Quote
    if quote, err = quoteService.FindRND(); err != nil {
        logger(err.Error())
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("что-то пошло не так %s", err.Error()))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    // make quoted image
    if err = drawer.AddQuoteBottom(quote, img); err != nil {
        logger(err.Error())
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("что-то пошло не так %s", err.Error()))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    // send result to chat
    buf := new(bytes.Buffer)
    if err := jpeg.Encode(buf, img, nil); err != nil {
        log.Printf("ERROR: %s", err.Error())
        return
    }

    msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Name: quote.Id, Bytes: buf.Bytes()})

    if _, err := bot.Send(msg); err != nil {
        log.Printf("Wednesday send error %s\n", err.Error())
    }
}
