package models

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/zktnotify/zktnotify/pkg/config"
	"github.com/zktnotify/zktnotify/pkg/xpath"
)

var (
	x      *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables,
		new(User),
		new(CardTime),
		new(Notify),
		new(Holiday),
	)
}

func NewEngine() {
	var (
		err       error
		params    = "?"
		driver    = ""
		sourceURL = ""
	)

	dbtype := config.Config.XServer.DB.Type
	if v := os.Getenv("ZKTNOTIFY_DBTYPE"); v != "" {
		dbtype = v
	}

	switch dbtype {
	case "sqlite3":
		driver = "sqlite3"
		sourceURL = config.Config.XServer.DB.Path

		if ok, _ := xpath.IsExists(sourceURL); !ok {
			os.MkdirAll(sourceURL, 0755)
		}

	case "mysql":
		driver = "mysql"
		sourceURL = fmt.Sprintf("%s:%s@tcp(%s)/%s%scharset=utf8mb4&parseTime=true",
			config.Config.XServer.DB.User,
			config.Config.XServer.DB.Password,
			config.Config.XServer.DB.Host,
			config.Config.XServer.DB.Name,
			params,
		)
	}

	x, err = xorm.NewEngine(driver, sourceURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, table := range tables {
		if ok, err := x.IsTableExist(table); !ok || err != nil {
			if err := x.Sync2(table); err != nil {
				log.Fatal(err)
			}
		}
	}
}
