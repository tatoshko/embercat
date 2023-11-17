package handlerTurbo

import (
    "encoding/json"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "math/rand"
    "net/url"
    "time"
)

func makeUserCollectionKey(key string, userID int64) string {
    return fmt.Sprintf(key, userID)
}

func GetRandomLiner() Liner {
    rand.Seed(time.Now().UnixNano())
    return NewLiner(rand.Intn(TOTAL_LINERS))
}

func GetChatMember(bot *tgbotapi.BotAPI, username string) (tgbotapi.ChatMember, error) {
    v := url.Values{}
    v.Add("chat_id", username)

    resp, err := bot.MakeRequest("getChatMember", v)
    if err != nil {
        return tgbotapi.ChatMember{}, err
    }

    var member tgbotapi.ChatMember
    err = json.Unmarshal(resp.Result, &member)

    return member, err
}
