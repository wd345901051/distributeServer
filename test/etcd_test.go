package test

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"testing"
	"time"
)

func TestEtcdPut(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.188.100:2379"},
		DialTimeout: 5 * time.Second,
	})
	defer client.Close()
	if err != nil {
		log.Fatal(err)
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	_, err = client.Put(timeout, "/client/1", "7878")
	if err != nil {
		log.Fatal(err)
	}
}

func TestEtcdGet(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.188.100:2379"},
		DialTimeout: 5 * time.Second,
	})
	defer client.Close()
	if err != nil {
		log.Fatal(err)
	}
	get, err := client.Get(context.TODO(), "/client", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range get.Kvs {
		fmt.Println(string(v.Value))
	}
}

func TestEtcdRM(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.188.100:2379"},
		DialTimeout: 5 * time.Second,
	})
	defer client.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Delete(context.TODO(), "/client", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
}
