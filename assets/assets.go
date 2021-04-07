package assets

import rice "github.com/GeertJohan/go.rice"

var box *rice.Box

func InitBox() {
    rice.MustFindBox(".")
}

func GetBox() *rice.Box {
    return box
}