package model

import "github.com/josexy/gw/global"

func MigrationTable() {
	_ = global.DB.AutoMigrate(
		&Admin{},
		&ServiceInfo{},
		&HttpRule{},
		&TcpRule{},
		&GrpcRule{},
		&LoadBalance{},
		&AccessControl{},
	)
}
