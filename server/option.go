package server

import (
	"github.com/wd345901051/distributeServer/internal/service"
	"time"
)

func WithKeepAliveTime(KeepAliveTime time.Duration) *service.Option {
	return &service.Option{F: func(o *service.Options) {
		o.KeepAliveTime = KeepAliveTime
	}}
}

func WithMaxKeepAliveTime(MaxKeepAliveTime time.Duration) *service.Option {
	return &service.Option{F: func(o *service.Options) {
		o.MaxKeepAliveTime = MaxKeepAliveTime
	}}
}

func WithMaxTryKeepAlive(MaxTryKeepAlive int) *service.Option {
	return &service.Option{F: func(o *service.Options) {
		o.MaxTryKeepAlive = MaxTryKeepAlive
	}}
}
