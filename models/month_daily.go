package models

import (
	"database/sql"
	"time"
)

const (
	SenderTimer = iota
	SenderUser
)

type MonthDaily struct {
	ID          int64     `xorm:"id"`
	UserID      uint64    `xorm:"user_id"`
	Month       int       `xorm:"month"`
	Sender      int8      `xorm:"sender"`
	CreatedTime time.Time `xorm:"created_time"`
	UpdatedTime time.Time `xorm:"updated_time"`
	Status      int8      `xorm:"status"`
}

func GetUserMonthDaily(uid uint64, month int) (*MonthDaily, error) {
	var m = MonthDaily{UserID: uid, Month: month}
	ok, err := x.Where("status=0").Limit(1, 0).Get(&m)
	if err != nil {
		return nil, err
	}
	if ok {
		return &m, nil
	}
	return nil, sql.ErrNoRows
}

func SaveUserMonthDaily(uid uint64, month int, sender int8) error {
	daily := MonthDaily{
		UserID:      uid,
		Month:       month,
		Sender:      sender,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	_, err := x.Insert(&daily)
	return err
}
