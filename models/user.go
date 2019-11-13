package models

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id            int64
	Name          string    `xorm:"name"`
	JobID         string    `xorm:"UNIQUE(UQE_USER) NOT NULL 'job_id'"`
	UserID        uint64    `xorm:"user_id"`
	Password      string    `xorm:"password"`
	NotifyURL     string    `xorm:"notify_url"`
	NotifyType    uint32    `xorm:"DEFAULT 0 'notify_type'"`
	NotifyAccount string    `xorm:"notify_account"`
	CreateTime    time.Time `xorm:"-"`
	CreateUnix    int64     `xorm:"'create_time'"`
	UpdateTime    time.Time `xorm:"-"`
	UpdateUnix    int64     `xorm:"'update_time'"`
	Status        int       `xorm:"DEFAULT 0 'status'"`
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
	rows, err := x.Rows(User{UserID: uid, Status: 0})
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

func GetUserByJobId(jobId string) *User {
	rows, err := x.Rows(User{JobID: jobId, Status: 0})
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

func SaveUser(user *User) error {
	affected, err := x.Insert(user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("save user failed")
	}
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

func ChangeUserStatus(jobId uint64,status int) error {
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
		Status: 0,
	}
	affected, err := x.Cols("user_id").Where("job_id=?", u.JobID).Update(&user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("not found job_id(%d) ", u.JobID)
	}
	return nil
}
