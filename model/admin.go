package model

import (
	"time"

	"github.com/josexy/gw/pkg/util"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"index;size:20;column:user_name;not null"`
	Password  string `gorm:"not null"`
	Salt      string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDelete  int `gorm:"column:is_delete"`
}

func (admin *Admin) TableName() string {
	return "gateway_admin"
}

func (admin *Admin) CheckPassword(password string) bool {
	saltPwd := util.GenerateSaltPassword(password, admin.Salt)
	return saltPwd == admin.Password
}

func (admin *Admin) Find(db *gorm.DB) error {
	err := db.Model(&Admin{}).Where("user_name = ?", admin.Username).First(admin).Error
	return err
}

func (admin *Admin) Update(db *gorm.DB) error {
	newPassword := admin.Password

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Admin{}).
			Where("user_name = ?", admin.Username).
			Find(admin).Error; err != nil {
			return err
		}

		admin.Password = util.GenerateSaltPassword(newPassword, admin.Salt)
		return tx.Model(&Admin{}).Where("id = ?", admin.ID).
			Update("password", admin.Password).Error
	})
	return err
}
