package core

import tba "github.com/go-telegram-bot-api/telegram-bot-api"

type Handler func(api *tba.BotAPI, update tba.Update)
