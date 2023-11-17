package handlerTurbo

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/draw"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
    "golang.org/x/image/math/fixed"
    "gopkg.in/redis.v3"
    "image"
    "image/color"
    "log"
    "math"
)

type Collection struct {
    connection *redis.Client
    userId     int64
    data       []redis.Z
    redisKey   string
}

func LoadCollection(redisClient *redis.Client, userId int64) (collection Collection, err error) {
    collection.userId = userId
    collection.redisKey = makeUserCollectionKey(REDIS_KEY_TURBO_COLLECTION, userId)
    collection.connection = redisClient
    var data []redis.Z
    if data, err = redisClient.ZRangeWithScores(collection.redisKey, 0, -1).Result(); err != nil {
        return Collection{}, errors.New("can't load collection")
    }

    collection.data = data

    log.Printf("DATA: %v, col.DATA %v. ERROR: %v", data, collection.data, err)

    if err != nil {
        return Collection{}, errors.New("empty collection")
    }

    return
}

func (c Collection) Count() int {
    return len(c.data)
}

func (c Collection) Has(liner Liner) bool {
    return c.ScoreOf(liner) > 0
}

func (c Collection) ScoreOf(liner Liner) int64 {
    log.Printf("DATA: %v", c.data)

    for _, v := range c.data {

        log.Printf("%s = %s ? %v", v.Member, liner.ID, v.Member == liner.ID)
        if v.Member == liner.ID {
            return int64(v.Score)
        }
    }

    return 0
}

func (c Collection) Add(liner Liner) (collection Collection, err error) {
    if _, err = c.connection.ZIncrBy(c.redisKey, 1, liner.ID).Result(); err != nil {
        return c, err
    }

    return LoadCollection(c.connection, c.userId)
}

// Removes one liner and returns new collection. If liner total count is one removes liner from collection
func (c Collection) RemoveOne(liner Liner) (collection Collection, err error) {
    if c.ScoreOf(liner) <= 1 {
        if _, err = c.connection.ZRem(c.redisKey, liner.ID).Result(); err != nil {
            return c, err
        }
    } else {
        if _, err = c.connection.ZIncrBy(c.redisKey, -1, liner.ID).Result(); err != nil {
            return c, err
        }
    }

    return LoadCollection(c.connection, c.userId)
}

// Moves liner from one collection to another
func (c Collection) MoveTo(collection Collection, liner Liner) (err error) {
    if !c.Has(liner) {
        return errors.New(fmt.Sprintf("В твоей коллекции нет вкладыша с номером <b>%s</b>", liner.ID))
    }

    if _, err = c.RemoveOne(liner); err != nil {
        return errors.New(fmt.Sprintf("Can't remove liner from collection because of %s", err.Error()))
    }

    if _, err = collection.Add(liner); err != nil {
        return errors.New(fmt.Sprintf("Can't add liner to collection because of of %s", err.Error()))
    }

    return
}

func (c Collection) GenerateCollectionPicture() *image.RGBA {
    // prepare canvas
    zeroPoint := image.Point{X: 0, Y: 0}
    canvasWH := image.Point{X: CANVAS_WIDTH, Y: CANVAS_HEIGHT}
    picWH := image.Point{X: PIC_WIDTH, Y: PIC_HEIGHT}

    collectionCanvas := image.NewRGBA(image.Rectangle{Min: zeroPoint, Max: canvasWH})
    resized := image.NewRGBA(image.Rectangle{Min: zeroPoint, Max: picWH})

    draw.Draw(collectionCanvas, collectionCanvas.Rect, image.White, zeroPoint, draw.Src)

    // prepare font
    ttf, _ := truetype.Parse(gobold.TTF)
    face := truetype.NewFace(ttf, &truetype.Options{
        Size:    28.0,
        DPI:     72.0,
        Hinting: font.HintingNone,
    })

    makePoint := func(i int) image.Point {
        y := int(math.Floor(float64(i/CANVAS_COLS)) * PIC_HEIGHT)
        x := (i % CANVAS_COLS) * PIC_WIDTH

        return image.Pt(x, y)
    }

    // Draw existing liners
    var err error
    var b []byte
    for i := 0; i < TOTAL_LINERS; i++ {
        liner, _ := NewLiner(i + 1)

        if c.Has(liner) {
            if b, err = liner.GetPicture(); err != nil {
                log.Printf("HandlerCollection liner.GetPicture error %s", err.Error())
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
            point := fixed.Point26_6{X: fixed.I(p.X + LABEL_OFFSET_X), Y: fixed.I(p.Y + LABEL_OFFSET_Y)}

            d := &font.Drawer{
                Dst:  collectionCanvas,
                Src:  image.NewUniform(color.RGBA{R: 80, G: 80, B: 80, A: 255}),
                Face: face,
                Dot:  point,
            }
            d.DrawString(liner.ID)
        }
    }

    draw.BiLinear.Scale(collectionCanvas, collectionCanvas.Rect, collectionCanvas, collectionCanvas.Bounds(), draw.Over, nil)

    return collectionCanvas
}
