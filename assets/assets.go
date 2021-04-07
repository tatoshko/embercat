package assets

import rice "github.com/GeertJohan/go.rice"

var box *rice.Box

func InitBox() {
    box = rice.MustFindBox(".")
}

func GetBox() *rice.Box {
    if box == nil {
        InitBox()
    }

    return box
}