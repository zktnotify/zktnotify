package models

import (
	"sort"
	"strings"
	"time"
)

var (
	EnableSQLite = false
)

type Times string

type TimeTag struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	UserID      uint64 `json:"userid"`
	Times       int    `json:"times"`
	CardDate    string `json:"card_date"`
	CardTimes   Times  `json:"card_times"`
	BadgeNumber string `json:"badgenumber"`
	DeptName    string `json:"DeptName"`
}

type TimeTagInfos struct {
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Rows  []TimeTag `json:"rows"`
}

func (t Times) Len() int {
	return len(strings.Split(strings.Trim(string(t), ","), ","))
}

func (t Times) Max() string {
	times := strings.Split(strings.Trim(string(t), ","), ",")
	sort.Strings(times)
	if len(times) == 0 {
		return ""
	}
	return times[len(times)-1]
}

func (t Times) Min() string {
	times := strings.Split(strings.Trim(string(t), ","), ",")
	sort.Strings(times)
	if len(times) == 0 {
		return ""
	}
	return times[0]
}

func (t Times) EveryTime() []string {
	times := strings.Split(strings.Trim(string(t), ","), ",")
	sort.Strings(times)
	return times
}

func (tag *TimeTagInfos) Split() map[string]TimeTag {
	data := make(map[string]TimeTag)
	for _, v := range tag.Rows {
		punch, _ := time.Parse("2006-01-02", v.CardDate)
		key := punch.Format("20060102")
		data[key] = v
	}
	return data
}

func (tag *TimeTagInfos) Today() *TimeTag {
	for _, r := range tag.Rows {
		if time.Now().Format("2006-01-02") == r.CardDate {
			return &TimeTag{
				ID:          r.ID,
				Name:        r.Name,
				UserID:      r.UserID,
				Times:       r.Times,
				CardDate:    r.CardDate,
				CardTimes:   r.CardTimes,
				BadgeNumber: r.BadgeNumber,
				DeptName:    r.DeptName,
			}
		}
	}
	return nil
}

func WorkingUsers() (users []User, err error) {
	dbSQL := `
		SELECT
			id,
			name,
			job_id,
			user_id,
			PASSWORD,
			notify_url,
			notify_type,
			notify_account,
			special_period_notify
		FROM
			user
		WHERE
			user_id NOT IN (
			SELECT
				user_id
			FROM
				notify
			WHERE
				card_date = ?
				AND status = 0
				AND (card_type = 2
				OR (card_type = 5 AND notified_status = 2))
			)
			AND status = 0
	`

	rows, err := x.DB().Query(dbSQL, time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		if err = rows.Scan(&u.Id, &u.Name, &u.JobID, &u.UserID, &u.Password, &u.NotifyToken, &u.NotifyType, &u.NotifyAccount, &u.SpecialPeriodNotify); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
