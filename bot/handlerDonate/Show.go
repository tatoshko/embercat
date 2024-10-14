package handlerDonate

import (
    "database/sql"
    "embercat/pgsql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func Show(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger := getLogger("SHOW")
    pg := pgsql.GetClient()

    q := `select username, sum(sum) as sum from donate group by username order by sum desc`

    var err error
    var rows *sql.Rows
    if rows, err = pg.Query(q); err != nil {
        logger(err.Error())
    }
    defer rows.Close()

    donates := NewDonates()
    for rows.Next() {
        donate := NewDonate()
        err = rows.Scan(&donate.Username, &donate.Sum)

        donates.Add(donate)
    }

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, getDonatesList(donates))
    msg.ParseMode = tgbotapi.ModeHTML
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Show error %s\n", err.Error())
    }
}
