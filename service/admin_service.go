package service

import (
	"errors"
	"strings"

	"github.com/josexy/gw/global"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/pkg/codes"
	"github.com/josexy/gw/serializer"
)

type AdminLoginService struct {
	UserName string `json:"username" binding:"required,min=0"`
	Password string `json:"password" binding:"required,min=0"`
}

func (service *AdminLoginService) Login() (serializer.Response, *model.Admin) {
	service.UserName = strings.TrimSpace(service.UserName)
	service.Password = strings.TrimSpace(service.Password)

	admin := model.Admin{
		Username: service.UserName,
		Password: service.Password,
	}

	if err := admin.Find(global.DB); err != nil {
		return serializer.BuildResponseErr(2000, err), nil
	}

	if !admin.CheckPassword(service.Password) {
		return serializer.BuildResponseErr(2001, errors.New("password validate failed")), nil
	}

	return serializer.BuildResponseOkWithData(codes.Success,
		serializer.Token{Token: admin.Username}), &admin
}

type AdminUpdateService struct {
	Password string `json:"password"`
}

func (service *AdminUpdateService) Update(username string) serializer.Response {

	admin := model.Admin{
		Username: username,
		Password: service.Password,
	}

	err := admin.Update(global.DB)

	if err != nil {
		return serializer.BuildResponseErr(2000, err)
	}
	return serializer.BuildResponseOk(codes.Success)
}
