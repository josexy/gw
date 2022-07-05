package global

type HttpProxyCommon struct {
	Addr string `yaml:"addr"`
}

type HttpProxyInfo struct {
	Port           int `yaml:"port"`
	ReadTimeout    int `yaml:"read_timeout"`
	WriteTimeout   int `yaml:"write_timeout"`
	MaxHeaderBytes int `yaml:"max_header_bytes"`
}
