package api

import (
	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/service"
)

func ServiceList(c *gin.Context) {
	var svc service.ListService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.ServiceList())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func ServiceDetail(c *gin.Context) {
	var svc service.DetailService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.ServiceDetail())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func ServiceDelete(c *gin.Context) {
	var svc service.DeleteService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.ServiceDelete())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func ServiceAddHTTP(c *gin.Context) {
	var svc service.AddHTTPService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.AddHTTPService())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func ServiceUpdateHTTP(c *gin.Context) {
	var svc service.UpdateHTTPService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.UpdateHTTPService())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func ServiceAddTcp(c *gin.Context) {
	var svc service.AddTCPService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.AddTCPService())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func ServiceUpdateTcp(c *gin.Context) {
	var svc service.UpdateTCPService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.UpdateTCPService())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func ServiceAddGrpc(c *gin.Context) {
	var svc service.AddGRPCService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.AddGRPCService())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}

func ServiceUpdateGrpc(c *gin.Context) {
	var svc service.UpdateGRPCService
	if err := c.ShouldBind(&svc); err == nil {
		ResponseJson(c, svc.UpdateGRPCService())
	} else {
		ResponseJsonErrorCode(c, 2000, err)
	}
}
