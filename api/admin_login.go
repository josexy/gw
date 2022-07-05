package api

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/pkg/codes"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/serializer"
	"github.com/josexy/gw/service"
)

func AdminLogin(c *gin.Context) {
	var svc service.AdminLoginService

	if err := c.ShouldBind(&svc); err == nil {
		resp, admin := svc.Login()

		if resp.Errno == codes.Success {
			sessionInfo := serializer.AdminSessionInfo{
				ID:        int(admin.ID),
				Username:  admin.Username,
				LoginTime: time.Now(),
			}
			data, err := json.Marshal(sessionInfo)
			if err != nil {
				ResponseJsonErrorCode(c, 2000, err)
				return
			}
			sess := sessions.Default(c)
			sess.Set(constants.AdminSessionInfoKey, string(data))
			sess.Save()
		}
		ResponseJson(c, resp)
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func AdminLogout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(constants.AdminSessionInfoKey)
	sess.Save()
	ResponseJson(c, serializer.BuildResponseOk(codes.Success))
}
