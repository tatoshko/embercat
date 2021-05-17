package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "regexp"
)

type CommandHandler func(api *tba.BotAPI, data string, chatID int64, update tba.Update)

type TBot struct {
    API        *tba.BotAPI
    commandMsg *regexp.Regexp
    updates    tba.UpdatesChannel

    HANDLERS map[string]CommandHandler
}

func Start(token, hook string) {
    bot := &TBot{
        API:        nil,
        commandMsg: regexp.MustCompile(`^/(?P<command>\w+)\s*(?P<data>.*)$`),
        HANDLERS:   make(map[string]CommandHandler),
    }

    bot.registerHandlers()

    if API, err := tba.NewBotAPI(token); err == nil {
        bot.API = API

        if _, err := API.SetWebhook(tba.NewWebhook(hook + "/" + token)); err != nil {
            log.Printf("SetHoook error %s\n", err.Error())
        }

        bot.updates = API.ListenForWebhook("/" + API.Token)
        bot.Watch()
    } else {
        log.Fatalf("NewAPIBot error %s\n", err.Error())
    }
}

func (bot *TBot) Watch() {
    for update := range bot.updates {
        if update.Message == nil {
            continue
        }

        text := update.Message.Text

        if bot.commandMsg.MatchString(text) {
            match := bot.parseCommand(bot.commandMsg, text)

            if handler, found := bot.HANDLERS[match["command"]]; found {
                log.Printf("Command: '%s', data: '%s'\n", match["command"], match["data"])
                go handler(bot.API, match["data"], update.Message.Chat.ID, update)
            }
        }
    }
}

func (bot *TBot) registerHandlers() {
    bot.HANDLERS["thread"] = handleThread
    bot.HANDLERS["day"] = handleWednesday
}

func (bot *TBot) parseCommand(r *regexp.Regexp, str string) (map[string]string) {
    match := r.FindStringSubmatch(str)
    subMatchMap := make(map[string]string)
    for i, name := range r.SubexpNames() {
        if i != 0 {
            subMatchMap[name] = match[i]
        }
    }

    return subMatchMap
}
