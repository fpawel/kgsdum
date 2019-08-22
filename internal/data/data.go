package data

import (
	"github.com/fpawel/gohelp"
	"github.com/fpawel/gohelp/winapp"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
)

//go:generate go run github.com/fpawel/gohelp/cmd/sqlstr/...

var (
	DB = func() *sqlx.DB {
		dir, err := winapp.AppDataFolderPath()
		if err != nil {
			panic(err)
		}
		fileName := filepath.Join(dir, "kgsdum", "kgsdum.sqlite")
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			panic("file " + fileName + " does not exist")
		}
		//fmt.Println(fileName)
		db := gohelp.OpenSqliteDBx(fileName)
		//db.MustExec(SQLCreate)
		return db
	}()
)
