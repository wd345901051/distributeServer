package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var Client *clientv3.Client

func init() {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.188.100:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	Client = c
}

func PutKVToEtcd(k, v string) error {
	_, err := Client.Put(context.TODO(), k, v)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetKVToEtcdByPrefix(prefix string) ([]string, error) {
	kvs, err := Client.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	vs := make([]string, len(kvs.Kvs))
	for _, c := range kvs.Kvs {
		vs = append(vs, string(c.Value))
	}
	return vs, nil
}

func DelKVToEtcd(k string) error {
	_, err := Client.Delete(context.TODO(), k)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
