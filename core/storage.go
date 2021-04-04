package core

import (
    "fmt"
    "github.com/boltdb/bolt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "strings"
)

var Storage *bolt.DB
var BUCKET_NAME = []byte("tbot-kv-storage")

func initStorage(db string) {
    if Storage, err = bolt.Open(db, 0600, nil); err == nil {
        defer Storage.Close()

        Storage.Update(func(tx *bolt.Tx) error {
            _, err := tx.CreateBucketIfNotExists(BUCKET_NAME)
            return err
        })
    } else {
        panic(err)
    }
}

func handleSet(bot *tba.BotAPI, update tba.Update, text string) {
    parts := strings.SplitAfterN(text, " ", 2)
    key, value := parts[0], parts[1]

    if Storage.IsReadOnly() {
        Storage.Update(func(tx *bolt.Tx) error {
            b := tx.Bucket(BUCKET_NAME)
            if err := b.Put([]byte(key), []byte(value)); err == nil {
                msg := tba.NewMessage(
                    update.Message.Chat.ID,
                    fmt.Sprintf("'%s' has been set to key '%s'", value, key),
                )
                msg.ReplyToMessageID = update.Message.MessageID

                bot.Send(msg)

                return nil
            } else {
                msg := tba.NewMessage(
                    update.Message.Chat.ID,
                    fmt.Sprintf("Error: %s", err.Error()),
                )
                bot.Send(msg)

                return err
            }
        })
    }
}

func handleGet(bot *tba.BotAPI, update tba.Update, key string) {
    if Storage.IsReadOnly() {
        Storage.View(func(tx *bolt.Tx) error {
            b := tx.Bucket(BUCKET_NAME)
            value := b.Get([]byte(key))

            msg := tba.NewMessage(update.Message.Chat.ID, string(value))
            msg.ReplyToMessageID = update.Message.MessageID

            bot.Send(msg)

            return nil
        })
    }
}
