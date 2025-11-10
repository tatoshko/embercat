package service

import (
    "database/sql"
    "embercat/pgsql"
)

type FrogReviewService struct {
    UserId int `json:"userId"`
}

func NewFrogReviewService(userId int) *FrogReviewService {
    return &FrogReviewService{UserId: userId}
}

func (s *FrogReviewService) Start() (err error) {
    pg := pgsql.GetClient()
    qSelect := `select 1 from frog_review where user_id = $1`

    row := pg.QueryRow(qSelect, s.UserId)
    if err = row.Scan(); err == sql.ErrNoRows {
        qInsert := `insert into frog_review (user_id) values ($1) on conflict do nothing`
        qFill := `insert into frog_review_item (select uuid_generate_v4(), $1, id, photoid from frog);`

        var tx *sql.Tx
        if tx, err = pg.Begin(); err != nil {
            return
        }
        defer tx.Rollback()

        tx.Exec(qInsert, s.UserId)
        tx.Exec(qFill, s.UserId)

        err = tx.Commit()
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
