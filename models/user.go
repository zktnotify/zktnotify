package models

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id                  int64
	Name                string    `xorm:"name"`
	JobID               string    `xorm:"UNIQUE(UQE_USER) NOT NULL 'job_id'"`
	UserID              uint64    `xorm:"user_id"`
	Password            string    `xorm:"password"`
	NotifyToken         string    `xorm:"notify_url"`
	NotifyType          uint32    `xorm:"DEFAULT 0 'notify_type'"` // 0 dingding, 1 serverChan, 2 wxpusher
	NotifyAccount       string    `xorm:"notify_account"`
	SpecialPeriodNotify bool      `xorm:"special_period_notify"` // 0 Not, 1 Yes
	CreateTime          time.Time `xorm:"-"`
	CreateUnix          int64     `xorm:"'create_time'"`
	UpdateTime          time.Time `xorm:"-"`
	UpdateUnix          int64     `xorm:"'update_time'"`
	Status              int       `xorm:"DEFAULT 0 UNIQUE(UQE_USER) NOT NULL 'status'"`
	NotifyCount         int       `xorm:"-"`
}

func (u *User) BeforeInsert() {
	u.CreateUnix = time.Now().Unix()
	u.UpdateUnix = time.Now().Unix()
}

func (u *User) AfterLoad() {
	u.CreateTime = time.Unix(u.CreateUnix, 0).Local()
	u.UpdateTime = time.Unix(u.UpdateUnix, 0).Local()
}

func (u *User) BeforeUpdate() {
	u.UpdateUnix = time.Now().Unix()
}

func GetUser(uid uint64) *User {
	rows, err := x.Where("status=0").Rows(User{UserID: uid})
	if err != nil {
		return nil
	}
	defer rows.Close()

	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user); err != nil {
			return nil
		}
		return &user
	}
	return nil
}

func GetUsers() []*User {
	rows, err := x.Where("status=0").Rows(&User{})
	if err != nil {
		return nil
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		if err := rows.Scan(user); err != nil {
			return nil
		}
		users = append(users, user)
	}
	return users
}

func GetUserByJobId(jobId string) *User {
	rows, err := x.Where("status=0").Rows(User{JobID: jobId})
	if err != nil {
		return nil
	}
	defer rows.Close()

	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user); err != nil {
			return nil
		}
		return &user
	}
	return nil
}

func IsTokenBind(token string) bool {
	cnt, err := x.Where("notify_url=? AND status = 0", token).Count(&User{})
	if err != nil {
		return false
	}
	return cnt > 0
}

func SaveUser(user *User) error {

	session := x.NewSession()
	session.Begin()
	defer session.Close()
	defer session.Rollback()

	affected, err := session.Insert(user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("save user failed")
	}
	session.Commit()
	return nil
}

func DeleteUser(jobId uint64) error {
	affected, err := x.Where("job_id=?", fmt.Sprintf("%d", jobId)).Delete(User{})
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("not found job_id(%d) ", jobId)
	}
	return nil
}

func ChangeUserStatus(jobId uint64, status int) error {
	user := User{
		Status: status,
	}
	affected, err := x.Cols("status").Where("job_id=?", fmt.Sprintf("%d", jobId)).Update(user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("not found job_id(%d) ", jobId)
	}
	return nil
}

func AllUsers() (users []User, err error) {
	rows, err := x.Where("status=0").Rows(User{})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *User) UpdateUserID() error {
	user := User{
		UserID: u.UserID,
	}
	affected, err := x.Cols("user_id").Where("job_id=? AND status=0", u.JobID).Update(&user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("not found job_id(%s) ", u.JobID)
	}
	return nil
}
