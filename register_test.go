package etcdkv

import (
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestNewRegister(t *testing.T) {

	register := NewRegister(
		RegisterClient(
			ClientEndpoints("127.0.0.1:2379,127.0.0.1:2389,127.0.0.1:2399"),
			ClientDialKeepAliveTime(time.Second*5),
			ClientDialTimeout(time.Second*5),
			ClientDialKeepAliveTimeout(time.Second*5),
		),
		RegisterTTL(time.Second*10),
		RegisterNamespace("/"),
		RegisterKvs("1", "1111:1:1:1"),
		RegisterKvs("2", "2222:2:2:2"),
		RegisterKvs("3", "3333:3:3:3"),
	)

	defer register.Close()
	register.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)
	<-c
}
