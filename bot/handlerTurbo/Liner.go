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
    ID       string
    box      *rice.Box
    filename string
}

func NewLiner(num int) Liner {
    return Liner{
        ID:       fmt.Sprintf("%03d", num),
        box:      assets.GetBox(),
        filename: fmt.Sprintf(TURBO_FILENAME_KEY, fmt.Sprintf("%d", num)),
    }
}

func NewLinerFromString(id string) (liner Liner, err error) {
    var num int64

    if num, err = strconv.ParseInt(id, 10, 32); err != nil {
        return Liner{}, errors.New("can't parse liner id from string")
    }

    return NewLiner(int(num)), nil
}

func (l Liner) GetPicture() (b []byte, err error) {
    if b, err = l.box.Bytes(l.filename); err != nil {
        log.Printf("HandlerTurbo box.Bytes error %s", err.Error())
        return nil, errors.New("can't load picture from box")
    }

    return b, nil
}
