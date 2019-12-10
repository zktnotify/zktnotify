// +build sqlite3

package models

import (
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	EnableSQLite = true
}
