package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strings"
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
	_, err = client.Put(timeout, "/client/1", "4545")
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

func TestEtcdWatch(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.188.100:2379"},
		DialTimeout: 5 * time.Second,
	})
	defer client.Close()
	if err != nil {
		log.Fatal(err)
	}
	watcher := clientv3.NewWatcher(client)
	watchChan := watcher.Watch(context.TODO(), "/client/", clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("发生了PUT操作！")
			case mvccpb.DELETE:
				fmt.Println("发生了DEL操作！")
			default:
				fmt.Println("Other Change")
			}
		}
	}
}

func TestGetIndex(t *testing.T) {
	s := "/client/10.10.109.51:8686"
	idx := strings.LastIndexByte(s, 'a')
	fmt.Println(idx)
}
