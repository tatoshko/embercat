package handlerPic

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "math/rand"
    "strings"
    "time"
)

const MAX = 409
const CDN = "https://pics.useful.team"

func Rnd(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    chatID := update.Message.Chat.ID

    rand.Seed(time.Now().UnixMicro())
    i := rand.Intn(MAX)
    id := fmt.Sprintf("%05d", i)

    msg := tgbotapi.NewMessage(chatID, strings.Join([]string{CDN, id}, "/"))

    if _, err := bot.Send(msg); err != nil {
        log.Printf("Pic.RND send message error: %s", err.Error())
    }
}
