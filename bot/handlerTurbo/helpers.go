package handlerTurbo

import (
    "encoding/json"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "gopkg.in/redis.v3"
    "net/url"
)

func makeLocalKey(key string, userID int) string {
    return fmt.Sprintf(key, userID)
}

func checkExists(stats []redis.Z, number string) bool {
    return GetScore(stats, number) > 0
}

func GetScore(stats []redis.Z, number string) float64 {
    for _, v := range stats {
        if v.Member == number {
            return v.Score
        }
    }

    return 0
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