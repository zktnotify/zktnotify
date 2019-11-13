package viewmodel

import (
	"fmt"
	"github.com/leaftree/ctnotify/models"
)

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name" validate:"required"`
	UserId   uint64 `json:"userId"`
	JobId    uint64 `json:"jobId" validate:"required"`
	Password string `json:"password" validate:"required"`
	Mobile   string `json:"mobile" validate:"required"`
	WebHook  string `json:"webHook"`
	Status   int    `json:"status"`
}

func (u *User) ToModelUser() *models.User {
	return &models.User{
		Id:            0,
		Name:          u.Name,
		JobID:         fmt.Sprintf("%d", u.JobId),
		UserID:        u.UserId,
		Password:      u.Password,
		NotifyURL:     u.WebHook,
		NotifyType:    0,
		NotifyAccount: u.Mobile,
		Status:        u.Status,
	}
}
