package client

import (
	"context"
	"fmt"
	"github.com/wd345901051/distributeServer/server/pb"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	Addr    string
	RPCConn *grpc.ClientConn
}

func NewClient(addr string, rpc *grpc.ClientConn) *Client {
	return &Client{
		Addr:    addr,
		RPCConn: rpc,
	}
}

func (c *Client) KeepAlive(maxTryKeepAlive int) error {
	rpcc := pb.NewServerClient(c.RPCConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := rpcc.KeepAlive(ctx, &pb.KeepAliveReq{})
	if err != nil {
		err = c.tryKeepAlive(maxTryKeepAlive)
		return err
	}
	fmt.Println(c.Addr, "yes")
	return nil
}
func (c *Client) tryKeepAlive(maxTryKeepAlive int) error {
	for i := 0; i < maxTryKeepAlive; i++ {
		time.Sleep(time.Millisecond)
		rpcc := pb.NewServerClient(c.RPCConn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		_, err := rpcc.KeepAlive(ctx, &pb.KeepAliveReq{})
		if err != nil {
			err = c.tryKeepAlive(maxTryKeepAlive)
			return err
		}
		cancel()
	}
	return nil
}
