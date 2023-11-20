package handlerDonate

type Donate struct {
    Username string
    Sum      float64
}

func NewDonate() *Donate {
    return &Donate{}
}

type Donates []*Donate

func NewDonates() *Donates {
    return &Donates{}
}

func (d Donates) Add(donate *Donate) {
    d = append(d, donate)
}
