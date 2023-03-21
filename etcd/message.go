package etcd

import "go.etcd.io/etcd/api/v3/mvccpb"

type Message struct {
	Type mvccpb.Event_EventType
	Msg  string
}

func newMessage(t mvccpb.Event_EventType, msg string) *Message {
	return &Message{
		Type: t,
		Msg:  msg,
	}
}
