package handlerTurbo

import (
    "embercat/assets"
    "errors"
    "fmt"
    rice "github.com/GeertJohan/go.rice"
    "log"
    "strconv"
)

type Liner struct {
    ID    int64
    Count int64
    box   *rice.Box
}

func NewLiner(num, count int64) (liner Liner, err error) {
    if num > TOTAL_LINERS || num <= 0 {
        return Liner{}, errors.New(fmt.Sprintf("There is existing liners only with numbers from 001 to %d", TOTAL_LINERS))
    }

    return Liner{
        ID:    num,
        Count: count,
        box:   assets.GetBox(),
    }, nil
}

func NewLinerFromString(id string) (liner Liner, err error) {
    var num int64
    if num, err = strconv.ParseInt(id, 10, 64); err != nil {
        return Liner{}, errors.New("can't parse liner id from string")
    }
    return NewLiner(num, 1)
}

func (l Liner) ToPicture() (b []byte, err error) {
    if b, err = l.box.Bytes(l.GetFilename()); err != nil {
        log.Printf("Roll box.Bytes error %s", err.Error())
        return nil, errors.New("can't load picture from box")
    }

    return b, nil
}

func (l Liner) GetFilename() string {
    return fmt.Sprintf(TURBO_FILENAME_KEY, l.ToString())
}

func (l Liner) ToString() string {
    return fmt.Sprintf("%03d", l.ID)
}
