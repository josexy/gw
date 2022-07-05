package load_balance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type UInt32Slice []uint32

func (s UInt32Slice) Len() int           { return len(s) }
func (s UInt32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s UInt32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type Hash func(data []byte) uint32

type ConsistentHashBalance struct {
	mux      sync.RWMutex
	hash     Hash
	replicas int
	keys     UInt32Slice
	hashMap  map[uint32]string
}

func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (c *ConsistentHashBalance) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("missing parameter")
	}
	addr := params[0]
	c.mux.Lock()
	defer c.mux.Unlock()
	// 虚拟节点
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr
	}
	// 对所有虚拟节点的哈希值进行排序，方便进行二分查找
	sort.Sort(c.keys)
	return nil
}

func (c *ConsistentHashBalance) Get(key string) (string, error) {
	hash := c.hash([]byte(key))

	idx := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })
	if idx == len(c.keys) {
		idx = 0
	}
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.hashMap[c.keys[idx]], nil
}

func (c *ConsistentHashBalance) Update(addrs []string) {
	c.keys = nil
	c.hashMap = map[uint32]string{}
	for _, addr := range addrs {
		c.Add(strings.Split(addr, ",")...)
	}
}
