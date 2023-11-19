package callbacks

import (
    "embercat/bot/core"
)

var (
    Callbacks = make(map[string]core.Handler)
)

func GetHandler(data string) (core.Handler, bool) {
    if handler, found := Callbacks[data]; found {
        return handler, true
    }

    return nil, false
}

func RegisterCallback(id string, f core.Handler) {
    Callbacks[id] = f
}

func UnregisterCallback(id string) {
    if _, found := Callbacks[id]; found {
        delete(Callbacks, id)
    }
}
