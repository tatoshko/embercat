package service

import "time"

type Place string

const (
    PlaceRNDMsg Place = "rnd_msg"
    PlaceENDPic Place = "rnd_pic"
)

type QuoteStat struct {
    id        string    `json:"id"`
    Which     Place     `json:"which"`
    CreatedAt time.Time `json:"created_at"`
}
