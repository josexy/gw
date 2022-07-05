package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/serializer"
)

func ResponseJsonErrorCode(c *gin.Context, code int, err error) {
	ResponseJson(c, serializer.BuildResponseErr(code, err))
}

func ResponseJson(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}
