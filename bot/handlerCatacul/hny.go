package handlerCatacul

import (
    "bytes"
    "embercat/assets"
    "fmt"
    tba "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
    "golang.org/x/image/math/fixed"
    "golang.org/x/text/feature/plural"
    "golang.org/x/text/language"
    "golang.org/x/text/message"
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "log"
    "time"
)

var (
    ttf     *truetype.Font
    face1   font.Face
    face2   font.Face
    printer *message.Printer
)

func Init() {
    ttf, _ = truetype.Parse(gobold.TTF)
    face1 = truetype.NewFace(ttf, &truetype.Options{
        Size:    42.0,
        DPI:     72.0,
        Hinting: font.HintingNone,
    })
    face2 = truetype.NewFace(ttf, &truetype.Options{
        Size:    56.0,
        DPI:     72.0,
        Hinting: font.HintingNone,
    })

    message.Set(language.Russian, "%.0f дней",
        plural.Selectf(1, "%.0f",
            "=0", "Усе, НГ!",
            "=1", "Один денек",
            "=2", "Два дня",
            "=3", "Три дня",
            "=4", "Четыре дня",
            "=5", "Пять дней",
            "=6", "Шесть дней",
            "=7", "Через семь дней",
            plural.One, "%.0f день",
            plural.Few, "%.0f дня",
            plural.Many, "%.0f дней",
            plural.Other, "%.0f дней",
        ),
    )

    printer = message.NewPrinter(language.Russian)
}

func Hny(bot *tba.BotAPI, update tba.Update) {
    pic := getSourcePic()

    currentTime := time.Now()
    newYear := Date(currentTime.Year()+1, 1, 1)

    days := newYear.Sub(currentTime).Hours() / 24

    addLabel(pic, 200, 100, fmt.Sprintf("До нового %d года", newYear.Year()), face1)
    addLabel(pic, 500, 156, printer.Sprintf("%.0f дней", days), face2)

    buf := new(bytes.Buffer)
    if err := jpeg.Encode(buf, pic, nil); err != nil {
        log.Printf("ERROR: %s", err.Error())
        return
    }

    msg := tba.NewPhotoUpload(update.Message.Chat.ID, tba.FileBytes{Name: "time", Bytes: buf.Bytes()})

    if _, err := bot.Send(msg); err != nil {
        log.Printf("Wednesday send error %s\n", err.Error())
    }
}

func addLabel(img *image.RGBA, x, y int, label string, face font.Face) {
    col := color.RGBA{255, 255, 255, 255}
    point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

    d := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(col),
        Face: face,
        Dot:  point,
    }
    d.DrawString(label)
}

func getSourcePic() (m *image.RGBA) {
    box := assets.GetBox()
    f, _ := box.Open("hny_cat.jpg")
    pic, _, _ := image.Decode(f)

    b := pic.Bounds()
    m = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
    draw.Draw(m, m.Bounds(), pic, b.Min, draw.Src)

    return
}

func Date(year, month, day int) time.Time {
    return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
