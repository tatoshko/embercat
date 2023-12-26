package handlerQuote

import (
    "embercat/bot/core"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "time"
)

type Quote struct {
    Id        string    `json:"id"`
    UserId    int       `json:"userId"`
    UserName  string    `json:"userName"`
    Text      string    `json:"text"`
    CreatedAt time.Time `json:"createdAt"`
}

func NewQuote() *Quote {
    return &Quote{}
}

func NewQuoteFromMessage(message *tgbotapi.Message) *Quote {
    q := NewQuote()
    q.UserId = message.From.ID
    q.UserName = core.GetUserName(message.From, true)
    q.Text = message.Text

    return q
}

func (q Quote) ToString() string {
    return fmt.Sprintf("%s. %s (c)", q.Text, q.UserName)
}
