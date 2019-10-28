package models

import (
	"database/sql"
	"time"
)

type Holiday struct {
	Id         int64
	Year       uint16    `xorm:"UNIQUE(UQE_HOLIDAY) NOT NULL 'year'"`
	Month      uint8     `xorm:"UNIQUE(UQE_HOLIDAY) NOT NULL 'month'"`
	Day        uint8     `xorm:"UNIQUE(UQE_HOLIDAY) NOT NULL 'day'"`
	WorkDay    bool      `xorm:"work_day"`
	CreateTime time.Time `xorm:"-"`
	CreateUnix int64     `xorm:"'create_time'"`
	UpdateTime time.Time `xorm:"-"`
	UpdateUnix int64     `xorm:"'update_time'"`
	Status     int       `xorm:"'status' default 0"`
}

func (h *Holiday) BeforeInsert() {
	h.CreateUnix = time.Now().Unix()
	h.UpdateUnix = time.Now().Unix()
}

func (h *Holiday) AfterLoad() {
	h.CreateTime = time.Unix(h.CreateUnix, 0).Local()
	h.UpdateTime = time.Unix(h.UpdateUnix, 0).Local()
}

func (h *Holiday) BeforeUpdate() {
	h.UpdateUnix = time.Now().Unix()
}

func GetHoliday(year uint16, month, day uint8) (*Holiday, error) {
	h := Holiday{
		Year:  year,
		Month: month,
		Day:   day,
	}

	ok, err := x.Where("status=0").Get(&h)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sql.ErrNoRows
	}

	return &h, nil
}
