package handlerTurbo

const REDIS_KEY_TURBO_COLLECTION = "turbo:%d:collection"
const REDIS_KEY_TURBO_DAY = "turbo:%d:day"
const TURBO_FILENAME_KEY = "turbo/2_turbo_old_%s.jpg"
const TOTAL_PICTURES = 330

const (
	PIC_WIDTH  = 500 / 3
	PIC_HEIGHT = 360 / 3
)

const (
	CANVAS_COLS   = 22
	CANVAS_ROWS   = 15
	CANVAS_WIDTH  = PIC_WIDTH * CANVAS_COLS
	CANVAS_HEIGHT = PIC_HEIGHT * CANVAS_ROWS
)

const (
	LABEL_OFFSET_X = (PIC_WIDTH / 2) - 30
	LABEL_OFFSET_Y = (PIC_HEIGHT / 2)
)

type CollectionItem struct {
	Number string
	Count  int
}
