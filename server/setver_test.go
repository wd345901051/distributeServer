package server

import (
	"github.com/wd345901051/distributeServer/pkg/serviceinfo"
	"log"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	info, err := serviceinfo.NewServiceInfo("10.10.109.51:8989", serviceinfo.MasterMethod)
	if err != nil {
		return
	}
	newServer, err := NewServer(info, WithKeepAliveTime(time.Second*2), WithMaxTryKeepAlive(3), WithMaxKeepAliveTime(time.Second*2))
	if err != nil {
		log.Fatal(err)
	}
	newServer.Run()
}

func TestNewClient(t *testing.T) {
	info, err := serviceinfo.NewServiceInfo("10.10.109.51:8686", serviceinfo.FollowerMethod)
	if err != nil {
		return
	}
	newServer, err := NewServer(info)
	if err != nil {
		log.Fatal(err)
	}
	newServer.Run()
}
