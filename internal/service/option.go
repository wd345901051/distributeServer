package service

import "time"

type Options struct {
	KeepAliveTime    time.Duration
	MaxKeepAliveTime time.Duration
	MaxTryKeepAlive  int
}

func (o *Options) Apply(opts []*Option) {
	for _, op := range opts {
		op.F(o)
	}
}

type Option struct {
	F func(o *Options)
}

func NewOptions(opts []*Option) *Options {
	o := &Options{KeepAliveTime: time.Second * 5, MaxTryKeepAlive: 3, MaxKeepAliveTime: 3}
	o.Apply(opts)
	return o
}
