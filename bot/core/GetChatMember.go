package core

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
)

func GetChatMember(bot *tba.BotAPI, chatID int64, userID int64) (member tba.ChatMember, err error) {
    config := tba.ChatConfigWithUser{
        ChatID: chatID,
        UserID: userID,
    }
    if member, err = bot.GetChatMember(tba.GetChatMemberConfig{ChatConfigWithUser: config}); err != nil {
        log.Printf("HandleReg GetChatMember err %s", err.Error())
    }

    return
}
