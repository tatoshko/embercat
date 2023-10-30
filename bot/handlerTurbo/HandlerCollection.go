package handlerTurbo

import (
	"bytes"
	"embercat/assets"
	redis2 "embercat/redis"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/redis.v3"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
)

func HandlerCollection(api *tgbotapi.BotAPI, update tgbotapi.Update) {
	redis := redis2.GetClient()
	if redis == nil {
		return
	}
	defer redis.Close()

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	localCollectionKey := makeLocalKey(REDIS_KEY_TURBO_COLLECTION, userID)

	if stats, err := redis.ZRangeWithScores(localCollectionKey, 0, -1).Result(); err != nil {
		log.Printf("HandlerCollection ZRangeWithScores error %s", err.Error())
	} else {
		message.Set(language.Russian, "В твоей коллекции %d вкладышей",
			plural.Selectf(1, "%d",
				"=0", "У тебя <b>нет</b> вкладышей",
				"=1", "У тебя пока <b>только один</b> вкладыш",
				"=2", "В твоей коллекции <b>Два</b> вкладыша",
				"=3", "В твоей коллекции <b>Три</b> вкладыша",
				"=4", "В твоей коллекции <b>Четыре</b> вкладыша",
				"=5", "В твоей коллекции <b>Пять</b> вкладышей",
				plural.One, "В твоей коллекции <b>%d</b> вкладыш",
				plural.Few, "В твоей коллекции <b>%d</b> вкладыша",
				plural.Many, "В твоей коллекции <b>%d</b> вкладышей",
			),
		)

		printer := message.NewPrinter(language.Russian)
		count := len(stats)

		result := printer.Sprintf("В твоей коллекции %d вкладышей", count)
		msg := tgbotapi.NewMessage(chatID, result)
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := api.Send(msg); err != nil {
			log.Printf("HandlerCollection send error %s", err.Error())
		}

		box := assets.GetBox()
		zeroPoint := image.Point{0, 0}
		collectionCanvas := image.NewRGBA(image.Rectangle{zeroPoint, image.Point{CANVAS_WIDTH, CANVAS_HEIGHT}})
		resized := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{PIC_WIDTH, PIC_HEIGHT}})
		draw.Draw(collectionCanvas, collectionCanvas.Rect, image.White, zeroPoint, draw.Src)

		ttf, _ := truetype.Parse(gobold.TTF)
		face := truetype.NewFace(ttf, &truetype.Options{
			Size:    28.0,
			DPI:     72.0,
			Hinting: font.HintingNone,
		})

		for i := 0; i < TOTAL_PICTURES; i++ {
			id := fmt.Sprintf("%03d", i+1)
			if checkExists(stats, id) {
				var b []byte
				var err error

				if b, err = box.Bytes(fmt.Sprintf(TURBO_FILENAME_KEY, id)); err != nil {
					log.Printf("HandlerCollection box.Bytes error %s", err.Error())
					continue
				}

				var pic image.Image

				if pic, _, err = image.Decode(bytes.NewReader(b)); err != nil {
					log.Printf("HandlerCollection image.Decode error %s", err.Error())
					continue
				}

				draw.BiLinear.Scale(resized, resized.Rect, pic, pic.Bounds(), draw.Src, nil)
				draw.Draw(collectionCanvas, resized.Bounds().Add(makePoint(i)), resized, zeroPoint, draw.Over)
			} else {
				p := makePoint(i)
				point := fixed.Point26_6{fixed.I(p.X + LABEL_OFFSET_X), fixed.I(p.Y + LABEL_OFFSET_Y)}

				d := &font.Drawer{
					Dst:  collectionCanvas,
					Src:  image.NewUniform(color.RGBA{80, 80, 80, 255}),
					Face: face,
					Dot:  point,
				}
				d.DrawString(id)
			}
		}

		draw.BiLinear.Scale(collectionCanvas, collectionCanvas.Rect, collectionCanvas, collectionCanvas.Bounds(), draw.Over, nil)

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, collectionCanvas, nil); err != nil {
			log.Printf("HandlerCollection jpeg.Encode error %s", err.Error())
		}

		msgp := tgbotapi.NewDocumentUpload(chatID, tgbotapi.FileBytes{Bytes: buf.Bytes(), Name: "Collection.jpg"})
		if _, err := api.Send(msgp); err != nil {
			log.Printf("HandlerCollection api.Send error %s", err.Error())
		}
	}
}

func makePoint(i int) image.Point {
	y := int(math.Floor(float64(i/CANVAS_COLS)) * PIC_HEIGHT)
	x := (i % CANVAS_COLS) * PIC_WIDTH

	return image.Pt(x, y)
}

func checkExists(stats []redis.Z, number string) bool {
	for _, v := range stats {
		if v.Member == number {
			return true
		}
	}

	return false
}
