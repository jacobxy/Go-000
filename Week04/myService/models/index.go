package models

import (
	pb "myService/api/newinfo"

	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type New struct {
	pb.NewInfo
}

func (n *New) TableName() string {
	return "news"
}

func init() {
	orm.RegisterModel(
		new(New),
	)
	var err error
	err = orm.RegisterDriver("sqlite3", orm.DRSqlite)
	err = orm.RegisterDataBase("default", "sqlite3", "data.db")
	err = orm.RunSyncdb("default", false, false)
	fmt.Println(err)
}
