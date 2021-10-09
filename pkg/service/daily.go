package service

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"sort"
	"time"

	"github.com/zktnotify/zktnotify/models"
	"github.com/zktnotify/zktnotify/pkg/config"
	"github.com/zktnotify/zktnotify/pkg/notify/typed"
	"github.com/zktnotify/zktnotify/pkg/tpl"
	"github.com/zktnotify/zktnotify/pkg/zkt"
)

func SendMonthDaily() {
	var (
		users = RetrieveAllUsers()
		month = int(time.Now().Month())
	)

	if len(config.Config.MonthDaily.Users) > 0 {
		user := []models.User{}
		for _, id := range config.Config.MonthDaily.Users {
			for _, val := range users {
				if val.UserID == uint64(id) {
					user = append(user, val)
					break
				}
			}
		}
		users = user
	}

	for _, user := range users {
		fmt.Println("send daily notify for", user.UserID)
		if user.NotifyType != uint32(typed.WXPUSHER) {
			continue
		}

		if _, err := models.GetUserMonthDaily(user.UserID, month); err != nil {
			if err != sql.ErrNoRows {
				log.Println("retrieve month daily failed:", err)
				continue
			}
		}

		msg := typed.Message{
			UID:        user.UserID,
			Name:       user.Name,
			Date:       time.Now().Format("2006-01-02"),
			Time:       time.Now().Format("15:04:05"),
			Status:     typed.MonthDaily,
			Token:      user.NotifyToken,
			Account:    user.NotifyAccount,
			NotifyType: typed.NotifierType(user.NotifyType),
		}
		if err := sendNotice(msg); err != nil {
			log.Printf("send take card notify for user(%d) failed:%v", user.UserID, err)
			continue
		}

		err := models.SaveUserMonthDaily(user.UserID, month, models.SenderTimer)
		if err != nil {
			log.Printf("save user(%d) month daily report failed:%v", user.UserID, err)
		}
	}
}

var dailytemp *template.Template

func init() {
	var err error
	const tplfile = "index.html"

	dailytemp, err = template.New("monthdaily").Parse(string(tpl.MustAsset(tplfile)))
	if err != nil {
		panic(err)
	}
}

func RetrieveMonthDaily(userid uint64, month int) (string, error) {
	if month == 0 {
		month = int(time.Now().Month())
	}
	user := models.GetUser(userid)
	if user == nil {
		return "", errors.New("user not found")
	}

	rpls, err := getUserMonthCardRecords(user, month)
	if err != nil {
		return "", err
	}

	var data = tpl.Reporter{
		Section: []tpl.ReportSection{{
			Line: rpls,
			Name: fmt.Sprintf("%d月份考勤报表", month),
		}},
	}

	var buff = bytes.NewBuffer(nil)
	if err := dailytemp.Execute(buff, data); err != nil {
		log.Println("gen report failed:", err)
		return "", err
	}
	return buff.String(), nil
}

func tomondate(t time.Time, day int) string {
	return time.Date(t.Year(), t.Month(), day, 1, 0, 0, 1, time.Local).Format("2006-01-02")
}
func lastmondate(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 1, 0, 0, 1, time.Local).Day()
}

func weekday(date string) string {
	cdText := map[time.Weekday]string{
		time.Sunday:    "星期日",
		time.Monday:    "星期一",
		time.Tuesday:   "星期二",
		time.Wednesday: "星期三",
		time.Thursday:  "星期四",
		time.Friday:    "星期五",
		time.Saturday:  "星期六",
	}

	t, _ := time.Parse("2006-01-02", date)
	return cdText[t.Weekday()]
}

func getUserMonthCardRecords(user *models.User, month int) ([]tpl.ReportLine, error) {
	var (
		rpl     []tpl.ReportLine
		thismon = time.Date(time.Now().Year(), time.Month(month), 1, 0, 0, 0, 1, time.Local)
		lastmon = thismon.AddDate(0, -1, 0)
		drange  = [][2]string{{tomondate(thismon, 1), tomondate(thismon, 20)}, {tomondate(lastmon, 21), tomondate(lastmon, lastmondate(lastmon))}}
	)

	tags, err := zkt.GetMonthTimeTag(user, drange)
	if err != nil {
		log.Printf("get user(%v) monthly card tags failed: %v", user.UserID, err)
	}

	if len(tags) > 0 {
		sort.Slice(tags, func(i, j int) bool { return tags[i].CardDate < tags[j].CardDate })
		for _, val := range tags {
			rpl = append(rpl, tpl.ReportLine{
				Date:     val.CardDate,
				Weekday:  weekday(val.CardDate),
				Times:    val.Times,
				Earliest: val.CardTimes.Min(),
				Latest:   val.CardTimes.Min(),
			})
		}
		return rpl, nil
	}

	records, err := models.GetUserCardDateRecords(user.UserID, tomondate(lastmon, 21), tomondate(thismon, 20))
	if err != nil {
		log.Printf("get user(%v) monthly card records failed: %v", user.UserID, err)
		return nil, err
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].CardDate == records[j].CardDate {
			return records[i].Times < records[j].Times
		}
		return records[i].CardDate < records[j].CardDate
	})

	mdata := map[string][]models.CardTime{}
	for _, val := range records {
		mdata[val.CardDate] = append(mdata[val.CardDate], val)
	}

	for date, tags := range mdata {
		min := tags[0]
		max := tags[len(tags)-1]
		rpl = append(rpl, tpl.ReportLine{
			Date:     date,
			Weekday:  weekday(date),
			Times:    int(max.Times),
			Earliest: min.CardTime,
			Latest:   max.CardTime,
		})
	}
	sort.Slice(rpl, func(i, j int) bool { return rpl[i].Date < rpl[j].Date })

	return rpl, nil
}
