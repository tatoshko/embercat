package core

import (
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "regexp"
)

var err error
var Bot *tba.BotAPI
var commandMsg = regexp.MustCompile(`^\/(?P<command>\w+)\s*(?P<data>.*)$`)

func InitBot(token, hook string) {
    if Bot, err = tba.NewBotAPI(token); err != nil {
        panic(err)
    }

    Bot.SetWebhook(tba.NewWebhook(hook + "/" + token))
    updates := Bot.ListenForWebhook("/" + Bot.Token)

    for update := range updates {
        text := update.Message.Text

        if commandMsg.MatchString(text) {
            match := reSubMatchMap(commandMsg, text)

            fmt.Printf("Command: %s, data: %s", match["command"], match["data"])

            switch match["command"] {
            case "set":
                handleSet(update, match["data"])
            case "get":
                handleGet(update, match["data"])
            }
        }
    }
}
