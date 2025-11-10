package service

type FrogReviewItem struct {
    Id           string `json:"id,omitempty"`
    FrogReviewId string `json:"frogReviewId,omitempty"`
    FrogId       string `json:"frogId,omitempty"`
    PhotoId      string `json:"photoId,omitempty"`
}

func NewFrogReviewItem() *FrogReviewItem {
    return &FrogReviewItem{}
}
