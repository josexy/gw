package global

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"net"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port       int    `yaml:"port"`
	Mode       string `yaml:"mode"`
	LogFileDir string `yaml:"logfile_dir"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type BaseConf struct {
	Server *ServerConfig `yaml:"server"`
	MySQL  *MySQLConfig  `yaml:"mysql"`
	Redis  *RedisConfig  `yaml:"redis"`
	Etcd   *EtcdConfig   `yaml:"etcd"`
}

type ProxyConf struct {
	Common *HttpProxyCommon `yaml:"common"`
	Http   *HttpProxyInfo   `yaml:"http"`
	Https  *HttpProxyInfo   `yaml:"https"`
}

var (
	AppConfig   = new(BaseConf)
	ProxyConfig = new(ProxyConf)
	DB          *gorm.DB
	Redis       redis.UniversalClient
	Logger      *logrus.Logger
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return ""
}

func InitConfig(baseConfPath, proxyConfPath string) {
	if config, err := ioutil.ReadFile(baseConfPath); err == nil {
		if err := yaml.Unmarshal(config, AppConfig); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	if config, err := ioutil.ReadFile(proxyConfPath); err == nil {
		if err := yaml.Unmarshal(config, ProxyConfig); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	ip := getLocalIP()
	if ip == "" {
		panic("not found local ip")
	}
	AppConfig.MySQL.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.MySQL.UserName, AppConfig.MySQL.Password,
		AppConfig.MySQL.Hostname, AppConfig.MySQL.Port,
		AppConfig.MySQL.DB)

	AppConfig.Redis.Addr = fmt.Sprintf("%s:%d", AppConfig.Redis.Hostname, AppConfig.Redis.Port)

	InitLogger(AppConfig.Server)
	InitMySQL(AppConfig.MySQL)
	InitRedis(AppConfig.Redis)
}
