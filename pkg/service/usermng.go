package service

import (
	"errors"
	"fmt"
	"github.com/leaftree/ctnotify/models"
	"github.com/leaftree/ctnotify/viewmodel"
	"strconv"
	"sync"
)

var (
	userMng     UserManager
	userMngOnce sync.Once
)

func GetUserManager() UserManager {
	userMngOnce.Do(func() {
		userMng = NewUserManager()
	})
	return userMng
}

type UserManager interface {
	AddUser(user *viewmodel.User) error
	GetUser(jobId uint64) (*viewmodel.User, error)
	GetUserAll() ([]*viewmodel.User, error)
	ChangeUserStatus(jobId uint64) error
	DeleteUser(jobId uint64) error
}

func NewUserManager() UserManager {
	return new(userManageImpl)
}

type userManageImpl struct {
	UserManager
}

func (*userManageImpl) AddUser(user *viewmodel.User) error {
	if u := models.GetUserByJobId(fmt.Sprintf("%d", user.JobId)); u != nil {
		return fmt.Errorf("The job id (%d) already exist!", user.JobId)
	}
	// add
	if err := models.SaveUser(user.ToModelUser()); err != nil {
		return err
	}
	return nil
}

func (*userManageImpl) GetUser(jobId uint64) (*viewmodel.User, error) {
	user := models.GetUserByJobId(fmt.Sprintf("%d", jobId))
	if user == nil {
		return nil, fmt.Errorf("The user (%d) not found!", jobId)
	}
	return &viewmodel.User{
		ID:       uint64(user.Id),
		Name:     user.Name,
		UserId:   user.UserID,
		JobId:    jobId,
		Password: user.Password,
		Mobile:   user.NotifyAccount,
		WebHook:  user.NotifyURL,
		Status:   user.Status,
	}, nil
}

func (*userManageImpl) DeleteUser(jobId uint64) error {
	return models.DeleteUser(jobId)
}

func (*userManageImpl) ChangeUserStatus(jobId uint64) error{
	user := models.GetUserByJobId(fmt.Sprintf("%d",jobId))
	if user == nil {
		return errors.New("The user not found!")
	}
	var status int
	if user.Status == 0 {
		status = 1
	} else {
		status = 0
	}
	return models.ChangeUserStatus(jobId,status)
}

func (*userManageImpl) GetUserAll() ([]*viewmodel.User, error) {
	users := models.GetUsers()
	if users == nil {
		return nil,errors.New("The user is empty")
	}
	_users := make([]*viewmodel.User,0)
	for _,item:=range users {
		jobId, _ := strconv.ParseInt(item.JobID, 10, 64)
		_users = append(_users, &viewmodel.User{
			ID:       uint64(item.Id),
			Name:     item.Name,
			UserId:   item.UserID,
			JobId:    uint64(jobId),
			Password: item.Password,
			Mobile:   item.NotifyAccount,
			WebHook:  item.NotifyURL,
			Status:   item.Status,
		})
	}
	return _users, nil
}
