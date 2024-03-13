package handlerTurbo

import (
    "bytes"
    "database/sql"
    "embercat/pgsql"
    "errors"
    "fmt"
    "github.com/golang/freetype/truetype"
    "golang.org/x/image/draw"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/gobold"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
    "log"
    "math"
)

type Collection struct {
    userId int64
    data   []Liner
}

func LoadCollection(userId int64) (collection Collection, err error) {
    pg := pgsql.GetClient()
    logger := getLogger("LoadCollection")

    collection.userId = userId
    q := `select linerid, count(linerid) from turbo where userid = $1 group by userid, linerid order by linerid`

    var rows *sql.Rows
    if rows, err = pg.Query(q, collection.userId); err != nil {
        logger(err.Error())
        return Collection{}, errors.New("unable to load collection")
    }

    for rows.Next() {
        var num, count int64
        if err = rows.Scan(&num, &count); err != nil {
            logger("SCAN ", err.Error())
            continue
        }

        var liner Liner
        if liner, err = NewLiner(num, count); err != nil {
            logger("LINER ", err.Error())
            continue
        }

        collection.data = append(collection.data, liner)
    }

    return
}

func (c Collection) GetDuplicates() (duplicates []Liner) {
    for _, liner := range c.data {
        if c.ScoreOf(liner) > 1 {
            duplicates = append(duplicates, liner)
        }
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
    for _, v := range c.data {
        if v.ID == liner.ID {
            return v.Count
        }
    }

    return 0
}

func (c Collection) Add(liner Liner) (collection Collection, err error) {
    pg := pgsql.GetClient()
    q := `insert into turbo (userid, linerid) values ($1, $2)`
    if _, err = pg.Exec(q, c.userId, liner.ID); err != nil {
        return
    }

    return LoadCollection(c.userId)
}

// RemoveOne liner and returns new collection. If liner total count is one removes liner from collection
func (c Collection) RemoveOne(liner Liner) (collection Collection, err error) {
    pg := pgsql.GetClient()
    q := `delete from turbo where createdat = (select createdat from turbo where userid = $1 and linerid = $2 order by createdat limit 1)`
    if _, err = pg.Exec(q, c.userId, liner.ID); err != nil {
        return
    }

    return LoadCollection(c.userId)
}

// MoveTo liner from one collection to another
func (c Collection) MoveTo(collection Collection, liner Liner) (err error) {
    if !c.Has(liner) {
        return errors.New(fmt.Sprintf("В твоей коллекции нет вкладыша с номером <b>%s</b>", liner.ToString()))
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
        liner, _ := NewLiner(int64(i+1), 1)

        if c.Has(liner) {
            if b, err = liner.ToPicture(); err != nil {
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
            d.DrawString(liner.ToString())
        }
    }

    draw.BiLinear.Scale(collectionCanvas, collectionCanvas.Rect, collectionCanvas, collectionCanvas.Bounds(), draw.Over, nil)

    return collectionCanvas
}
