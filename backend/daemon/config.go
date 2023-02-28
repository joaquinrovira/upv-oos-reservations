package daemon

import (
	"fmt"

	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
)

type Listen string

const (
	Socket Listen = "socket"
	Net    Listen = "net"
)

type DaemonConfig struct {
	Listen Listen
	Socket struct {
		Filename string
	}
	Net struct {
		Interface string
		Port      uint16
	}
	Test         map[int]int
	AnotherTest  map[string]string
	Test2        [10]int
	AnotherTest2 []int
}

func (cfg *DaemonConfig) Populate() {
	cfg.setDefaults()
	err := util.SetFromEnv(cfg, "DaemonConfig.Config")
	fmt.Println(err)
	fmt.Println(cfg)
}

func (cfg *DaemonConfig) setDefaults() {
	cfg.Listen = Socket
	cfg.Socket.Filename = "daemon.sock"
	cfg.Net.Interface = ""
	cfg.Net.Port = 8080
	cfg.Test = map[int]int{}
	cfg.AnotherTest = map[string]string{}
	cfg.AnotherTest2 = []int{1, 2, 3}
}
