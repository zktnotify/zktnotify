package models

import (
	"time"
)

type CardTime struct {
	Id          int64
	UserID      uint64    `xorm:"user_id"`
	Times       uint64    `xorm:"times"`
	CardDate    string    `xorm:"card_date"`
	CardTime    string    `xorm:"card_time"`
	BadgeNumber string    `xorm:"badge_number"`
	CreateTime  time.Time `xorm:"-"`
	CreateUnix  int64     `xorm:"'create_time'"`
	UpdateTime  time.Time `xorm:"-"`
	UpdateUnix  int64     `xorm:"'update_time'"`
	Status      int       `xorm:"status"`
}

func (c *CardTime) BeforeInsert() {
	c.CreateUnix = time.Now().Unix()
	c.UpdateUnix = time.Now().Unix()
}

func (c *CardTime) AfterLoad() {
	c.CreateTime = time.Unix(c.CreateUnix, 0).Local()
	c.UpdateTime = time.Unix(c.UpdateUnix, 0).Local()
}

func (c *CardTime) BeforeUpdate() {
	c.UpdateUnix = time.Now().Unix()
}

func CardTimes(c CardTime) (times []CardTime, err error) {
	rows, err := x.Where("status=0").Rows(c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ct := CardTime{}
		if err = rows.Scan(&ct); err != nil {
			return nil, err
		}
		times = append(times, ct)
	}

	return times, nil
}

func (c CardTime) Punched() error {
	session := x.NewSession()
	session.Begin()
	defer session.Close()

	if n, err := session.Insert(&c); err != nil || n <= 0 {
		session.Rollback()
		return err
	}
	return session.Commit()
}

func EarliestAndLatestCardTime(uid uint64, cdate string) (*CardTime, *CardTime) {
	var times []CardTime

	if x.Where("user_id=? AND card_date=? AND status=0", uid, cdate).Asc("card_time").Find(&times) != nil {
		return nil, nil
	}
	if len(times) == 0 {
		return nil, nil
	}

	return &times[0], &times[len(times)-1]
}

func GetUserCardDateRecords(uid uint64, from, to string) ([]CardTime, error) {
	var times []CardTime

	err := x.Where("user_id=? AND card_date BETWEEN ? AND ? AND status=0", uid, from, to).Asc("card_time").Find(&times)
	return times, err
}
