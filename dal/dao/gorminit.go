package dao

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/kackerx/go-mall/config"
)

var (
	_DbMaster *gorm.DB
	_DbSlave  *gorm.DB
)

func DB() *gorm.DB {
	return _DbSlave
}

func DBMaster() *gorm.DB {
	return _DbMaster
}

func init() {
	_DbMaster = initDB(config.Conf.DB.Master)
	_DbSlave = initDB(config.Conf.DB.Slave)
}

func initDB(option *config.DBConnectOption) *gorm.DB {
	db, err := gorm.Open(mysql.Open(option.Dsn), &gorm.Config{
		Logger:                                   NewGormLogger(500 * time.Millisecond),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxOpenConns(option.MaxOpen)
	sqlDb.SetConnMaxLifetime(time.Duration(option.MaxLiftTime) * time.Second * 60)
	sqlDb.SetMaxIdleConns(option.MaxIdle)

	if err = sqlDb.Ping(); err != nil {
		panic(err)
	}

	return db
}
