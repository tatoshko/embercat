package handlerDonate

import (
    "database/sql"
    "embercat/pgsql"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func Show(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger := getLogger("SHOW")
    pg := pgsql.GetClient()

    q := `select username, sum from donate order by sum desc`

    var err error
    var rows *sql.Rows
    if rows, err = pg.Query(q); err != nil {
        logger(err.Error())
    }
    defer rows.Close()

    donates := NewDonates()
    for rows.Next() {
        var username string
        var sum float64
        err = rows.Scan(&username, &sum)
        donates.Add(username, sum)
    }

    logger(fmt.Sprintf("%q", donates))

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, getDonatesList(donates))
    msg.ParseMode = tgbotapi.ModeHTML
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Show error %s\n", err.Error())
    }
}
