package handlerPic

import (
    "fmt"
    "math/rand"
    "strings"
    "time"
)

const MAX = 409
const CDN = "https://pics.useful.team"

type Pic struct {
    link string
}

func NewPic(n int) *Pic {
    id := fmt.Sprintf("%05d", n)
    return &Pic{link: strings.Join([]string{CDN, id}, "/")}
}

func NewPicRND() *Pic {
    rand.Seed(time.Now().UnixMicro())
    return NewPic(rand.Intn(MAX))
}

func (p Pic) URL() string {
    return p.link
}
