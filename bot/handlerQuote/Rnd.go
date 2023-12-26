package handlerQuote

import (
    "bytes"
    "embercat/bot/handlerQuote/drawer"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "image"
    "image/draw"
    "image/jpeg"
    "log"
    "net/http"
)

func Make(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("RND")
    chatID := update.Message.Chat.ID

    replay := update.Message.ReplyToMessage

    if replay == nil || replay.Photo == nil {
        msg := tgbotapi.NewMessage(chatID, "Ты че пёс, сообщение должно быть с картинкой")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    var photoID string
    for _, pic := range *replay.Photo {
        photoID = pic.FileID
    }

    var fileURL string
    if fileURL, err = bot.GetFileDirectURL(photoID); err != nil {
        msg := tgbotapi.NewMessage(chatID, "Не смог получить картинку")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    img, err := getSourceImg(fileURL)

    service := NewService()

    var quote *Quote
    if quote, err = service.findRND(); err != nil {
        logger(err.Error())
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("что-то пошло не так %s", err.Error()))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    drawer.AddLabel(img, quote.ToString())

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

func getSourceImg(fileURL string) (m *image.RGBA, err error) {
    var resp *http.Response
    if resp, err = http.Get(fileURL); err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var pic image.Image
    pic, _, err = image.Decode(resp.Body)

    b := pic.Bounds()
    m = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
    draw.Draw(m, m.Bounds(), pic, b.Min, draw.Src)

    return
}
