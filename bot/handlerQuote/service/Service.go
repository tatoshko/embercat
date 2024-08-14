package service

import (
    "embercat/pgsql"
    "errors"
    "strings"
)

var (
    ErrEmptyText = errors.New("nothing to save")
)

type Service struct {
}

func NewService() *Service {
    return &Service{}
}

func (s Service) Add(quote *Quote) (err error) {
    if strings.Trim(quote.Text, " ") == "" {
        return ErrEmptyText
    }

    pg := pgsql.GetClient()
    q := `insert into quote (chat_id, user_id, username, text) values ($1, $2, $3, $4)`

    _, err = pg.Exec(q, quote.ChatID, quote.UserId, quote.UserName, quote.Text)

    return
}

func (s Service) FindRND(chatId int64) (quote *Quote, err error) {
    pg := pgsql.GetClient()
    q := `select id, user_id, username, text, created_at from quote where chat_id = $1 order by random() limit 1`

    row := pg.QueryRow(q, chatId)

    if row.Err() != nil {
        return nil, row.Err()
    }

    quote = NewQuote()
    row.Scan(
        &quote.Id,
        &quote.ChatID,
        &quote.UserId,
        &quote.UserName,
        &quote.Text,
        &quote.CreatedAt,
    )

    return
}
