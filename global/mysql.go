package global

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	DSN      string
}

func InitMySQL(config *MySQLConfig) {
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		Logger.Errorln(err)
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		Logger.Errorln(err)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	db.Set("gorm:table_options", "ENGINE=Innodb DEFAULT CHARSET=utf8mb4")
	DB = db
}
