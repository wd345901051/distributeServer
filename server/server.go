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
	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	opt       *service.Options
	svcInfo   *serviceinfo.ServiceInfo
	Followers map[string]*client.Client
	CEtcd     chan *etcd.Message
	pb.UnimplementedServerServer
}

func NewServer(svcInfo *serviceinfo.ServiceInfo, opts ...*service.Option) (*Server, error) {
	if svcInfo == nil {
		return nil, errors.New("")
	}
	return &Server{
		opt:       service.NewOptions(opts),
		svcInfo:   svcInfo,
		Followers: make(map[string]*client.Client, 0),
		CEtcd:     make(chan *etcd.Message),
	}, nil
}

func (s *Server) Run() {
	fmt.Println("[ Server ] Server Is Running By "+s.svcInfo.Address+" , And This Method Is ", s.svcInfo.EtcdPrefix[1:len(s.svcInfo.EtcdPrefix)-1], ", Success ! ")
	go s.RegisterRPC()
	go s.registerEtcdServer()
	go s.getFollowersByEtcd()
	go s.TODO()
	s.Stop()
}

func (s *Server) TODO() {
	for {
		time.Sleep(time.Second * 10)
		select {}
	}
}

func (s *Server) KeepAlive(context.Context, *pb.KeepAliveReq) (*pb.KeepAliveRsp, error) {
	fmt.Println("心跳检测被调用了!")
	return &pb.KeepAliveRsp{}, nil
}

func (s *Server) RegisterRPC() {
	listen, err := net.Listen("tcp", s.svcInfo.Address)
	if err != nil {
		return
	}
	gs := grpc.NewServer()
	pb.RegisterServerServer(gs, s)
	fmt.Println(s.svcInfo.EtcdPrefix[1:len(s.svcInfo.EtcdPrefix)-1] + " RegisterRPC Success!")
	if err := gs.Serve(listen); err != nil {
		return
	}
}

func (s *Server) RegisterRPCFollowers(address []string) {
	for _, addr := range address {
		s.AddRPCFollower(addr)
	}
	go s.WatchEtcdChan()
	go s.KeepAliveToFs()
}

func (s *Server) Stop() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig
	fmt.Println("GC start!")
	for k, _ := range s.Followers {
		s.RemoveRPCFollower(k)
	}
	s.cancellationEtcdServer()
	etcd.Client.Close()
	fmt.Println("GC end!")
	os.Exit(0)
}

// KeepAliveToFs RegisterRPCFollowers() Success After
func (s *Server) KeepAliveToFs() {
	for {
		time.Sleep(s.opt.KeepAliveTime)
		if len(s.Followers) == 0 {
			fmt.Println("No Followers Is Online Status !")
			continue
		}
		fmt.Println("发送心跳检测!")
		for _, f := range s.Followers {
			err := f.KeepAlive(s.opt.MaxTryKeepAlive, s.opt.MaxKeepAliveTime)
			if err != nil {
				fmt.Println("KeepAlive Error To ", f.Addr, err)
				break
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
	_, err := etcd.Client.Delete(context.TODO(), s.svcInfo.EtcdPrefix+s.svcInfo.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Server) getFollowersByEtcd() {
	if s.svcInfo.Method != 1 {
		return
	}
	address, version, err := etcd.GetKVToEtcdByPrefix("/client/")
	if err != nil {
		fmt.Println(err)
		return
	}
	go etcd.WatchKV("/client/", version, s.CEtcd)
	s.RegisterRPCFollowers(address)
}

func (s *Server) AddRPCFollower(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	c := client.NewClient(addr, conn)
	s.Followers[addr] = c
}

func (s Server) WatchEtcdChan() {
	for msg := range s.CEtcd {
		switch msg.Type {
		case mvccpb.PUT:
			fmt.Println("客户端", msg.Msg, "上线了！")
			s.AddRPCFollower(msg.Msg)
		case mvccpb.DELETE:
			fmt.Println("客户端", msg.Msg, "下线了！")
			s.RemoveRPCFollower(msg.Msg)
		default:
			fmt.Println("Other Change")
		}
	}
}

func (s *Server) RemoveRPCFollower(addr string) {
	fmt.Println(addr)
	fmt.Println(s.Followers)
	if conn, ok := s.Followers[addr]; ok {
		conn.RPCConn.Close()
		delete(s.Followers, addr)
	}
}
