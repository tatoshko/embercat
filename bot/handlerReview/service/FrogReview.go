package service

import "time"

type FrogReview struct {
    Id        string    `json:"id,omitempty"`
    UserId    int64     `json:"userId,omitempty"`
    StartedAt time.Time `json:"startedAt"`
}

func NewFrogReview() *FrogReview {
    return &FrogReview{}
}
