package load_balance

import (
	"errors"
	"hash/crc32"
	"strings"
)

type IPHashBalance struct {
	rss []string
}

func (ih *IPHashBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("missing parameter")
	}
	addr := params[0]
	ih.rss = append(ih.rss, addr)
	return nil
}

func (ih *IPHashBalance) Get(key string) (string, error) {
	if len(ih.rss) == 0 {
		return "", nil
	}
	curAddr := ih.rss[crc32.ChecksumIEEE([]byte(key))%uint32(len(ih.rss))]
	return curAddr, nil
}

func (ih *IPHashBalance) Update(addrs []string) {
	ih.rss = ih.rss[:0]
	for _, addr := range addrs {
		ih.Add(strings.Split(addr, ",")...)
	}
}
