package model

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/ulricqin/blog/config"
)

var Orm *xorm.Engine

func InitMysql() {
	var err error
	Orm, err = xorm.NewEngine("mysql", config.G.Mysql.Addr)
	if err != nil {
		log.Fatalln("fail to connect mysql", err)
	}

	Orm.SetMaxIdleConns(config.G.Mysql.Idle)
	Orm.SetMaxOpenConns(config.G.Mysql.Max)
	Orm.ShowSQL(config.G.Mysql.ShowSql)
	Orm.Logger().SetLevel(core.LOG_INFO)
}
