package core

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "strconv"
)

func GetChatMember(bot *tgbotapi.BotAPI, chatID int64, userID string) (member tgbotapi.ChatMember, err error) {
    if id, err := strconv.Atoi(userID); err == nil {
        config := tgbotapi.ChatConfigWithUser{
            ChatID: chatID,
            UserID: id,
        }
        if member, err = bot.GetChatMember(config); err != nil {
            log.Printf("HandleReg GetChatMember err %s", err.Error())
        }
    } else {
        log.Printf("USER_ID '%s', %q", err.Error(), id)
    }

    return
}
