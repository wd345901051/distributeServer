package serviceinfo

import "errors"

type ServiceInfo struct {
	Address    string
	Method     int
	EtcdPrefix string
}

const (
	MasterMethod   = 1
	FollowerMethod = 2
)

func NewServiceInfo(addr string, method int) (*ServiceInfo, error) {
	etcdPrefix := ""
	switch method {
	case 1:
		etcdPrefix = "/server/"
	case 2:
		etcdPrefix = "/client/"
	default:
		return nil, errors.New("server Method Is Invalid ")
	}
	return &ServiceInfo{
		Address:    addr,
		Method:     method,
		EtcdPrefix: etcdPrefix,
	}, nil
}
