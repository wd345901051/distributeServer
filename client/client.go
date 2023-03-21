package client

import (
	"context"
	"errors"
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

func (c *Client) KeepAlive(maxTryKeepAlive int, MaxKeepAliveTime time.Duration) error {
	if maxTryKeepAlive == 0 {
		return errors.New("Out Of MaxTryKeepAlive !")
	}
	rpcc := pb.NewServerClient(c.RPCConn)
	ctx, cancel := context.WithTimeout(context.Background(), MaxKeepAliveTime)
	defer cancel()
	_, err := rpcc.KeepAlive(ctx, &pb.KeepAliveReq{})
	if err != nil {
		// err try
		return c.KeepAlive(maxTryKeepAlive-1, MaxKeepAliveTime)
	}
	fmt.Println(c.Addr, "yes")
	return nil
}
