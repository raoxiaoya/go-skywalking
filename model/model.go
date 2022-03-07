package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	conf2 "github.com/phprao/go-skywalking.git/conf"
	"github.com/phprao/go-skywalking.git/tracerhelper/gormagent"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var DbDB *sql.DB
var yamlFile = "../conf/app.yaml"

// Setup initializes the database instance
func Setup() {
	conf, err := conf2.ReadYamlConfig(yamlFile)
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		conf.Apps.Database.Username,
		conf.Apps.Database.Password,
		conf.Apps.Database.Host,
		conf.Apps.Database.Port,
		conf.Apps.Database.Database)
	Db, DbDB = ConnectDatabase(dsn, 10, 500)
	_ = Db.Use(gormagent.SetGormPlugin("db"))
}

// CloseDB closes database connection (unnecessary)
func CloseAllDb() {
	_ = DbDB.Close()
}

func ConnectDatabase(dsn string, maxIdleConns int, maxOpenConns int) (*gorm.DB, *sql.DB) {
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("gormDB.Setup err: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("gormDB.Setup err: %v", err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Minute * 3)

	return gormDB, sqlDB
}
