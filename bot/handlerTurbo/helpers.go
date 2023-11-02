package handlerTurbo

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "gopkg.in/redis.v3"
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

func GetChatMember(bot *tgbotapi.BotAPI, chatID int64, username string) (tgbotapi.Chat, error) {
    config := tgbotapi.ChatConfig{
        ChatID:             chatID,
        SuperGroupUsername: username,
    }

    return bot.GetChat(config)
}
