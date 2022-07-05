package serializer

import (
	"time"
)

type AdminSessionInfo struct {
	ID        int       `json:"id"`
	Username  string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
}

type AdminInfo struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

type Token struct {
	Token string `json:"token"`
}
