package etcd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/josexy/gw/pkg/constants"
)

type customLoadBalanceObserver struct{}

func (obs *customLoadBalanceObserver) Update(addrs []string) {
	fmt.Printf("Observer update: %v\n", addrs)
}

func TestEtcdGetPut(t *testing.T) {

	const prefix = constants.EtcdPrefix
	endpoints := []string{"127.0.0.1:2379"}
	observer := new(customLoadBalanceObserver)

	discovery, err := NewServiceDiscovery(endpoints, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	discovery.AddObserver(observer)

	t.Logf("watch server: %v", prefix)
	_, err = discovery.WatchService(prefix)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		service, err := NewServiceRegister(endpoints, 5)
		if err != nil {
			t.Fatal(err)
		}

		hostPort1 := net.JoinHostPort("127.0.0.1", "1001")
		hostPort2 := net.JoinHostPort("127.0.0.1", "1002")

		key1 := prefix + hostPort1
		key2 := prefix + hostPort2
		t.Log(key1)
		t.Log(key2)

		_ = service.RegisterService(key1, hostPort1)
		_ = service.RegisterService(key2, hostPort2)

		time.Sleep(time.Second * 6)
		t.Log("service register close")
		service.Close()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)

	<-interrupt
}
