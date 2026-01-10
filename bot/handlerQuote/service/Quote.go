package service

import (
    "embercat/bot/core"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "regexp"
    "time"
)

type Quote struct {
    Id        string    `json:"id"`
    ChatID    int64     `json:"chatID"`
    UserId    int64     `json:"userId"`
    UserName  string    `json:"userName"`
    Text      string    `json:"text"`
    CreatedAt time.Time `json:"createdAt"`
}

func NewQuote() *Quote {
    return &Quote{}
}

func NewQuoteFromMessage(message *tgbotapi.Message) *Quote {
    q := NewQuote()
    q.ChatID = message.Chat.ID
    q.UserId = message.From.ID
    q.UserName = core.GetUserName(message.From, false)
    q.Text = message.Text

    return q
}

func (q *Quote) Len() int {
    return len(q.Text)
}

func (q *Quote) Words() []string {
    r := regexp.MustCompile(`\S+`)
    return r.FindAllString(q.Text, -1)
}

func (q *Quote) ToString() string {
    return fmt.Sprintf("%s\n%s (c)", q.Text, q.UserName)
}
