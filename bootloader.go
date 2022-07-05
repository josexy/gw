package main

import (
	"errors"
	"flag"
)

type Bootloader struct {
	dashboard     bool
	server        bool
	baseConfPath  string
	proxyConfPath string
}

func NewBootloader() *Bootloader {
	loader := &Bootloader{}
	flag.BoolVar(&loader.dashboard, "dashboard", false, "start dashboard")
	flag.BoolVar(&loader.server, "server", false, "start server")
	flag.StringVar(&loader.baseConfPath, "base_conf", "conf/config.yaml", "base config file path")
	flag.StringVar(&loader.proxyConfPath, "proxy_conf", "conf/proxy.yaml", "proxy config file path")
	flag.Parse()
	return loader
}

func (loader *Bootloader) IsDashboard() bool {
	return loader.dashboard
}

func (loader *Bootloader) IsServer() bool {
	return loader.server
}

func (loader *Bootloader) Check() error {
	if loader.server && loader.dashboard {
		return errors.New("dashboard or server")
	}
	switch {
	case loader.dashboard:
		return nil
	case loader.server:
		return nil
	default:
		return errors.New("you must start dashboard or server")
	}
}
