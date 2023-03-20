package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/wd345901051/distributeServer/client"
	"github.com/wd345901051/distributeServer/etcd"
	"github.com/wd345901051/distributeServer/internal/service"
	"github.com/wd345901051/distributeServer/pkg/serviceinfo"
	"github.com/wd345901051/distributeServer/server/pb"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

type Server struct {
	opt       *service.Options
	svcInfo   *serviceinfo.ServiceInfo
	Followers map[string]*client.Client
	pb.UnimplementedServerServer
}

func NewServer(svcInfo *serviceinfo.ServiceInfo, opts ...*service.Option) (*Server, error) {
	if svcInfo == nil {
		return nil, errors.New("")
	}
	s := &Server{}
	s.svcInfo = svcInfo
	s.opt = service.NewOptions(opts)
	return s, nil
}

func (s *Server) Run() {
	go s.RegisterRPC()
	go s.registerEtcdServer()
	s.getFollowersByEtcd()
	sig := make(chan os.Signal)
	<-sig
	s.Stop()
}

func (s *Server) KeepAlive(context.Context, *pb.KeepAliveReq) (*pb.KeepAliveRsp, error) {
	return &pb.KeepAliveRsp{}, nil
}

func (s *Server) RegisterRPC() {
	listen, err := net.Listen("tcp", s.svcInfo.Address)
	if err != nil {
		return
	}
	gs := grpc.NewServer()
	pb.RegisterServerServer(gs, s)
	fmt.Println("RegisterRPC Success!")
	if err := gs.Serve(listen); err != nil {
		return
	}
}

func (s *Server) RegisterRPCFollowers(addrs []string) {
	for _, addr := range addrs {
		conn, err := grpc.Dial(addr, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}
		c := client.NewClient(addr, conn)
		s.Followers[addr] = c
	}
	go s.KeepAliveToFs()
}

func (s *Server) Stop() {
	fmt.Println("GC start!")
	for _, f := range s.Followers {
		f.RPCConn.Close()
	}
	s.cancellationEtcdServer()
	etcd.Client.Close()
	os.Exit(0)
}

func (s *Server) KeepAliveToFs() {
	for {
		time.Sleep(s.opt.KeepAliveTime)
		for _, f := range s.Followers {
			err := f.KeepAlive(s.opt.MaxTryKeepAlive)
			if err != nil {
				fmt.Println("KeepAlive Error To ", f.Addr, err)
			}
		}
	}
}

func (s Server) registerEtcdServer() {
	err := etcd.PutKVToEtcd(s.svcInfo.EtcdPrefix+s.svcInfo.Address, s.svcInfo.Address)
	if err != nil {
		return
	}
}

func (s *Server) cancellationEtcdServer() {
	_, err := etcd.Client.Delete(context.TODO(), "/server/"+s.svcInfo.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Server) getFollowersByEtcd() {
	if s.svcInfo.Method != 1 {
		return
	}
	address, err := etcd.GetKVToEtcdByPrefix("/client/")
	if err != nil {
		fmt.Println(err)
		return
	}
	s.RegisterRPCFollowers(address)
}
