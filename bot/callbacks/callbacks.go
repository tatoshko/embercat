package callbacks

import "embercat/bot/types"

var (
	Callbacks = make(map[string]types.Handler)
)

func GetHandler(data string) (types.Handler, bool) {
	if handler, found := Callbacks[data]; found {
		return handler, true
	}

	return nil, false
}

func RegisterCallback(id string, f types.Handler) {
	Callbacks[id] = f
}

func UnregisterCallback(id string) {
	if _, found := Callbacks[id]; found {
		delete(Callbacks, id)
	}
}
