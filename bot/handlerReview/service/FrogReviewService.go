package service

import (
    "database/sql"
    "embercat/pgsql"
)

type FrogReviewService struct {
    UserId int64 `json:"userId"`
}

func NewFrogReviewService(userId int64) *FrogReviewService {
    return &FrogReviewService{UserId: userId}
}

func (s *FrogReviewService) FindById(itemId string) (reviewItem *FrogReviewItem, err error) {
    pg := pgsql.GetClient()
    q := `select id, frog_review_id, frog_id, photo_id from frog_review_item where id = $1`

    reviewItem = NewFrogReviewItem()
    row := pg.QueryRow(q, itemId)
    err = row.Scan(
        &reviewItem.Id,
        &reviewItem.FrogReviewId,
        &reviewItem.FrogId,
        &reviewItem.PhotoId,
    )

    return
}

func (s *FrogReviewService) Start() (id string, err error) {
    pg := pgsql.GetClient()
    qSelect := `select id from frog_review where user_id = $1`

    row := pg.QueryRow(qSelect, s.UserId)
    if err = row.Scan(&id); err == sql.ErrNoRows {
        qInsert := `insert into frog_review (user_id) values ($1) on conflict do nothing`

        if _, err = pg.Exec(qInsert, s.UserId); err != nil {
            return
        }

        row := pg.QueryRow(qSelect, s.UserId)
        if err = row.Scan(&id); err != nil {
            return
        }

        qFill := `insert into frog_review_item (select uuid_generate_v4(), $1, id, photoid from frog);`
        _, err = pg.Exec(qFill, id)
    }

    return
}

func (s *FrogReviewService) Stop() (err error) {
    pg := pgsql.GetClient()
    q := `delete from frog_review where user_id = $1`

    _, err = pg.Exec(q, s.UserId)

    return
}

func (s *FrogReviewService) Next() (reviewItem *FrogReviewItem, err error) {
    pg := pgsql.GetClient()
    q := `select id, frog_review_id, frog_id, photo_id from frog_review_item order by id desc limit 1`

    row := pg.QueryRow(q)

    reviewItem = NewFrogReviewItem()
    err = row.Scan(
        &reviewItem.Id,
        &reviewItem.FrogReviewId,
        &reviewItem.FrogId,
        &reviewItem.PhotoId,
    )

    return
}

func (s *FrogReviewService) Approve(item *FrogReviewItem) (err error) {
    pg := pgsql.GetClient()
    q := `delete from frog_review_item where id = $1`

    _, err = pg.Exec(q, item.Id)

    return
}

func (s *FrogReviewService) Reject(item *FrogReviewItem) (err error) {
    if err = s.Approve(item); err != nil {
        return
    }

    pg := pgsql.GetClient()
    q := `delete from frog where id = $1`

    _, err = pg.Exec(q, item.FrogId)

    return
}
