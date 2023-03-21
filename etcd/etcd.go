package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
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

func GetKVToEtcdByPrefix(prefix string) ([]string, int64, error) {
	kvs, err := Client.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	vs := make([]string, len(kvs.Kvs))
	for _, c := range kvs.Kvs {
		vs = append(vs, string(c.Value))
	}
	watchStartRevision := kvs.Header.Revision + 1
	return vs, watchStartRevision, nil
}

func DelKVToEtcd(k string) error {
	_, err := Client.Delete(context.TODO(), k)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func WatchKV(k string, version int64, c chan *Message) {
	watcher := clientv3.NewWatcher(Client)
	watchChan := watcher.Watch(context.TODO(), k, clientv3.WithPrefix(), clientv3.WithRev(version))
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			fmt.Println(event)
			switch event.Type {
			case mvccpb.PUT:
				c <- newMessage(mvccpb.PUT, string(event.Kv.Value))
			case mvccpb.DELETE:
				s := string(event.Kv.Key)
				idx := strings.LastIndexByte(s, '/')
				if idx == -1 {
					continue
				}
				c <- newMessage(mvccpb.DELETE, s[idx+1:])
			}
		}
	}
}
