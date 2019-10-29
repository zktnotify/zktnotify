package models

import (
	"database/sql"
	"time"
)

type Notify struct {
	Id             int64
	UserID         uint64    `xorm:"UNIQUE(UQE_NOTIFY) NOT NULL 'user_id'"`
	CardDate       string    `xorm:"UNIQUE(UQE_NOTIFY) NOT NULL 'card_date'"`
	CardTime       string    `xorm:"UNIQUE(UQE_NOTIFY) NOT NULL 'card_time'"`
	CardType       uint64    `xorm:"UNIQUE(UQE_NOTIFY) NOT NULL 'card_type'"`
	NotifiedStatus uint64    `xorm:"'notified_status'"`
	CreateTime     time.Time `xorm:"-"`
	CreateUnix     int64     `xorm:"'create_time'"`
	UpdateTime     time.Time `xorm:"-"`
	UpdateUnix     int64     `xorm:"'update_time'"`
	Status         int       `xorm:"'status'"`
}

func (n *Notify) BeforeInsert() {
	n.CreateUnix = time.Now().Unix()
	n.UpdateUnix = time.Now().Unix()
}

func (n *Notify) AfterLoad() {
	n.CreateTime = time.Unix(n.CreateUnix, 0).Local()
	n.UpdateTime = time.Unix(n.UpdateUnix, 0).Local()
}

func (n *Notify) BeforeUpdate() {
	n.UpdateUnix = time.Now().Unix()
}

func IsNotified(uid uint64, cdate string, ctype uint64) bool {
	n := Notify{UserID: uid, CardDate: cdate, CardType: ctype, NotifiedStatus: 1}
	if ok, err := x.Where("status=0").Exist(&n); !ok || err != nil {
		return false
	}
	return true
}

func Notified(uid uint64, ctype uint64, cdate, ctime string) error {
	n := Notify{
		UserID:         uid,
		CardDate:       cdate,
		CardTime:       ctime,
		CardType:       ctype,
		NotifiedStatus: 1,
	}

	ok, err := x.Where("status=0").Exist(&n)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	_, err = x.Insert(&n)
	return err
}

func CanNotify(uid uint64, cdate string) bool {
	n := Notify{UserID: uid, CardDate: cdate, CardType: 5}
	if ok, _ := x.Where("notified_status=2 AND status=0").Exist(&n); ok {
		return false
	}
	return true
}

func UpdateNotice(uid uint64, ctype uint64, cdate, ctime string) error {
	n := Notify{
		UserID:   uid,
		CardDate: cdate,
		CardTime: ctime,
		CardType: ctype,
	}

	ok, err := x.Where("notified_status=1 AND status=0").Exist(&n)
	if err != nil {
		return err
	}
	if ok {
		_, err := x.Update(&n)
		return err
	}

	_, err = x.Insert(&n)
	return err
}

func CounterNotice(uid uint64, cdate string) error {
	n := Notify{
		NotifiedStatus: 2,
	}
	affected, err := x.Where("user_id=? AND card_date=? AND card_type=? AND status=0", uid, cdate, 5).Update(&n)
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
