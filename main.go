package main

import (
	"github.com/josexy/gw/global"
	"github.com/josexy/gw/logx"
	"github.com/josexy/gw/model"
	"github.com/josexy/gw/router"
)

// make build && ./goserver -dashboard
// make build && ./goserver -server
func main() {
	loader := NewBootloader()
	if err := loader.Check(); err != nil {
		logx.Error("bootloader failed: %v", err)
		return
	}
	global.InitConfig(loader.baseConfPath, loader.proxyConfPath)
	model.MigrationTable()
	if loader.IsDashboard() {
		logx.Debug("start dashboard")
		svr := NewServer(global.AppConfig.Server.Port, router.NewRouter())
		svr.Run()
	} else if loader.IsServer() {
		logx.Debug("start server")
		svr := ProxyServer{}
		svr.Run()
	}
	logx.Debug("bye :)")
}
