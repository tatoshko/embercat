package handlerPic

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "time"
)

func Rnd(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    chatID := update.Message.Chat.ID

    pic := NewPicRND()

    //resp, err := http.Get(pic.URL())
    //if err != nil {
    //    logger("HTTP.GET", err.Error())
    //    return
    //}
    //defer resp.Body.Close()
    //
    //if resp.StatusCode != 200 {
    //    logger("STATUS", string(resp.StatusCode))
    //    return
    //}
    //
    //img, _, err := image.Decode(resp.Body)

    msg := tgbotapi.NewPhotoUpload(chatID, tgbotapi.File{FilePath: pic.URL()})
    msg.Caption = fmt.Sprintf("%v", time.Now())

    if _, err := bot.Send(msg); err != nil {
        log.Printf("Pic.RND send message error: %s", err.Error())
    }
}
