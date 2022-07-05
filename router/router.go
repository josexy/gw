package router

import (
	"io/ioutil"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/josexy/gw/api"
	"github.com/josexy/gw/global"
	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/middleware"
	"github.com/josexy/gw/validator"
)

func NewRouter() *gin.Engine {

	if global.AppConfig.Server.Mode == "release" {
		logx.DebugMode = false
		gin.DefaultWriter = ioutil.Discard
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	validator.InitValidator()

	store, err := sessions.NewRedisStore(10, "tcp", global.AppConfig.Redis.Addr,
		global.AppConfig.Redis.Password, []byte("secret"))
	if err != nil {
		logx.Error("sessions new redis store: %v", err)
	}
	r.Use(middleware.RateLimit(), middleware.Cors(), middleware.Logger())

	adminLogin := r.Group("/admin_login")
	adminLogin.Use(
		sessions.Sessions("mysession", store),
	)
	{
		adminLogin.POST("login", api.AdminLogin)
		adminLogin.GET("logout", api.AdminLogout)
	}

	admin := r.Group("/admin")
	admin.Use(
		sessions.Sessions("mysession", store),
		middleware.SessionAuthMiddleware(),
	)
	{
		admin.GET("/admin_info", api.AdminInfo)
		admin.POST("/change_pwd", api.AdminChangePwd)
	}

	service := r.Group("/service")
	service.Use(
		sessions.Sessions("mysession", store),
		middleware.SessionAuthMiddleware(),
	)
	{
		service.GET("/service_list", api.ServiceList)
		service.GET("/service_delete", api.ServiceDelete)
		service.GET("/service_detail", api.ServiceDetail)
		service.POST("/service_add_http", api.ServiceAddHTTP)
		service.POST("/service_update_http", api.ServiceUpdateHTTP)
		service.POST("/service_add_tcp", api.ServiceAddTcp)
		service.POST("/service_update_tcp", api.ServiceUpdateTcp)
		service.POST("/service_add_grpc", api.ServiceAddGrpc)
		service.POST("/service_update_grpc", api.ServiceUpdateGrpc)
	}
	return r
}
