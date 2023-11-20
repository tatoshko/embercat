package handlerDonate

import (
    "embercat/pgsql"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
)

func Add(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
    userID := update.Message.From.UserName

    if userID != COLLECTOR {
        msg.Text = getUserDonateMessage()
    } else {
        if donater, sum, err := parseArgs(update.Message.CommandArguments()); err != nil {
            msg.Text = fmt.Sprintf("Ошибонька: %s", err.Error())
        } else {
            if err := addDonater(donater, sum); err != nil {
                log.Printf("HandlerAdd error %s\n", err.Error())
                return
            } else {
                msg.Text = getDonateMessage(donater, sum)
            }
        }
    }

    if _, err := bot.Send(msg); err != nil {
        log.Printf("HandlerAdd error %s\n", err.Error())
    }
}

func addDonater(donater string, sum float64) (err error) {
    pg := pgsql.GetClient()

    q := `insert into donate (username, sum) values ($1, $2)`
    _, err = pg.Exec(q, donater, sum)

    return
}
