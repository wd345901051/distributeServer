package test

import (
	"context"
	"fmt"
	"github.com/wd345901051/distributeServer/pkg/serviceinfo"
	"github.com/wd345901051/distributeServer/server"
	"github.com/wd345901051/distributeServer/server/pb"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	info, err := serviceinfo.NewServiceInfo("127.0.0.1:8989", serviceinfo.MasterMethod)
	if err != nil {
		return
	}
	newServer, err := server.NewServer(info)
	if err != nil {
		log.Fatal(err)
	}
	newServer.Run()
}

func TestNewClient(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8989", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req, err := c.KeepAlive(ctx, &pb.KeepAliveReq{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req.String(), "yes")
}
