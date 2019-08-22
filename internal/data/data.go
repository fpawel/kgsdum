package data

import (
	"github.com/fpawel/gohelp"
	"github.com/jmoiron/sqlx"
	"github.com/lxn/win"
	"os"
	"path/filepath"
	"syscall"
)

//go:generate go run github.com/fpawel/gohelp/cmd/sqlstr/...

var (
	DB = func() *sqlx.DB {

		var buf [win.MAX_PATH]uint16
		if !win.SHGetSpecialFolderPath(0, &buf[0], win.CSIDL_APPDATA, false) {
			panic("SHGetSpecialFolderPath failed")
		}
		dir := syscall.UTF16ToString(buf[0:])

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
