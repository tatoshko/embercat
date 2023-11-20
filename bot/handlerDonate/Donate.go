package handlerDonate

type Donate struct {
    Username string
    Sum      float64
}

func NewDonate() *Donate {
    return &Donate{}
}

type Donates struct {
    data []*Donate
}

func NewDonates() *Donates {
    return &Donates{}
}

func (d *Donates) Add(donate *Donate) {
    d.data = append(d.data, donate)
}

func (d Donates) GetAll() []*Donate {
    return d.data
}
