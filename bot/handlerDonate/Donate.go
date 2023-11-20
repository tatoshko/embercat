package handlerDonate

type Donate struct {
    Username string
    Sum      float64
}

func NewDonate(username string, sum float64) *Donate {
    return &Donate{Username: username, Sum: sum}
}

type Donates []*Donate

func NewDonates() *Donates {
    return &Donates{}
}

func (d Donates) Add(username string, sum float64) {
    d = append(d, NewDonate(username, sum))
}
