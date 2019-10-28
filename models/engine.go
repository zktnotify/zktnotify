package models

import (
	"log"
	"os"
	"path"

	"github.com/go-xorm/xorm"
	"github.com/leaftree/onoffice/pkg/config"
	"github.com/leaftree/onoffice/pkg/xpath"
	_ "github.com/mattn/go-sqlite3"
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
	var err error
	dbfile := config.Config.XServer.File.DB
	dbpath := path.Dir(config.Config.XServer.File.DB)

	if ok, _ := xpath.IsExists(dbpath); !ok {
		os.MkdirAll(dbpath, 0755)
	}

	x, err = xorm.NewEngine("sqlite3", dbfile)
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
