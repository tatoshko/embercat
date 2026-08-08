package handlerWednesday

import (
    "database/sql"
    "embercat/pgsql"
    "log"
)

type Frog struct {
    Id      string
    PhotoId string
    PHash   string
}

type FrogService struct {
    Threshold float64
}

func NewFrogService() *FrogService {
    return &FrogService{Threshold: 10.0}
}

func (f *FrogService) Add(photoId string, vector string) (err error) {
    pg := pgsql.GetClient()
    q := `insert into frog (photoId, phash) VALUES ($1, $2::vector)`

    _, err = pg.Exec(q, photoId, vector)

    return
}

func (f *FrogService) IsExists(vector string) (isExists bool, frog *Frog, err error) {
    pg := pgsql.GetClient()
    q := `select id, photoid, phash, phash <-> $1::vector as distance from frog order by phash <-> $1::vector limit 1`

    frog = &Frog{}
    var distance float64
    err = pg.QueryRow(q, vector).Scan(&frog.Id, &frog.PhotoId, &frog.PHash, &distance)

    if err == sql.ErrNoRows {
        return false, nil, nil
    }

    if err != nil {
        return false, nil, err
    }

    log.Printf("Distance %.f", distance)

    return distance < f.Threshold, frog, nil
}
