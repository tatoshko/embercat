package ai

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandlerTextAnalizer(API *tgbotapi.BotAPI, update tgbotapi.Update) {
    //var err error
    //
    //redisInst := redisServ.GetClient()
    //if redisInst == nil {
    //    return
    //}
    //defer redisInst.Close()
    //
    //if update.Message.Text == "" {
    //    return
    //}
    //
    //symbols := strings.Split(update.Message.Text, "")
    //
    //if _, err = redisInst.SAdd(REDIS_KEY_SYMBOL, symbols...).Result(); err != nil {
    //    log.Printf("Sadd error %s", err.Error())
    //}
    //
    //r := regexp.MustCompile(`(\w+|[А-я])`)
    //parts := r.FindAllString(update.Message.Text, -1)
    //
    //var hash map[string]int
    //for _, word := range parts {
    //
    //}

}
