package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/pkg/codes"
	"github.com/josexy/gw/pkg/constants"
	"github.com/josexy/gw/serializer"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if adminInfo, ok := session.Get(constants.AdminSessionInfoKey).(string); !ok || adminInfo == "" {
			c.JSON(http.StatusOK, serializer.BuildResponseErr(codes.InternalErrorCode, errors.New("user not login")))
			c.Abort()
			return
		}
		c.Next()
	}
}
