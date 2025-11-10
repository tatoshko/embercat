package service

import "time"

type Place string

const (
    PlaceRNDMsg Place = "rnd_msg"
    PlaceENDPic Place = "rnd_pic"
)

type QuoteStat struct {
    id        string    `json:"id"`
    ChatId    int64     `json:"chatId"`
    Which     Place     `json:"which"`
    CreatedAt time.Time `json:"created_at"`
}
