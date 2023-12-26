package handlerQuote

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
    q := `insert into quote (userid, username, text) values ($1, $2, $3)`

    _, err = pg.Exec(q, quote.UserId, quote.UserName, quote.Text)

    return
}

func (s Service) findRND() (quote *Quote, err error) {
    pg := pgsql.GetClient()
    q := `select id, userId, userName, text, createdAt from quote order by random() limit 1`

    row := pg.QueryRow(q)

    if row.Err() != nil {
        return nil, row.Err()
    }

    quote = NewQuote()
    row.Scan(
        &quote.Id,
        &quote.UserId,
        &quote.UserName,
        &quote.Text,
        &quote.CreatedAt,
    )

    return
}
