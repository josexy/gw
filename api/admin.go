package api

import (
	"encoding/json"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/pkg/codes"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/serializer"
	"github.com/josexy/gw/service"
)

func AdminInfo(c *gin.Context) {
	var adminSessInfo serializer.AdminSessionInfo

	sess := sessions.Default(c)
	sessInfo := sess.Get(constants.AdminSessionInfoKey)
	if err := json.Unmarshal([]byte(sessInfo.(string)), &adminSessInfo); err != nil {
		ResponseJson(c, serializer.BuildResponseErr(2000, err))
		return
	}
	adminInfo := serializer.AdminInfo{
		ID:           adminSessInfo.ID,
		Name:         adminSessInfo.Username,
		LoginTime:    adminSessInfo.LoginTime,
		Avatar:       "https://i0.hdslb.com/bfs/article/b44b901909fa721b1fbb5d5145f5cfccb22abd45.jpg",
		Introduction: "I am an administrator",
		Roles:        []string{"admin"},
	}
	ResponseJson(c, serializer.BuildResponseOkWithData(codes.Success, adminInfo))
}

func AdminChangePwd(c *gin.Context) {
	var adminSessInfo serializer.AdminSessionInfo

	sess := sessions.Default(c)
	sessInfo := sess.Get(constants.AdminSessionInfoKey)
	if err := json.Unmarshal([]byte(sessInfo.(string)), &adminSessInfo); err != nil {
		ResponseJson(c, serializer.BuildResponseErr(2000, err))
		return
	}
	var svc service.AdminUpdateService
	if err := c.ShouldBind(&svc); err == nil {
		resp := svc.Update(adminSessInfo.Username)
		ResponseJson(c, resp)
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}
