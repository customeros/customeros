package util

import (
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"google.golang.org/grpc"
)

type DialFactory interface {
	GetOasisAPICon() (*grpc.ClientConn, error)
}

type DialFactoryImpl struct {
	conf *config.Config
}

func (dfi DialFactoryImpl) GetOasisAPICon() (*grpc.ClientConn, error) {
	return grpc.Dial(dfi.conf.Service.OasisApiUrl, grpc.WithInsecure())
}

func MakeDialFactory(conf *config.Config) DialFactory {
	dfi := new(DialFactoryImpl)
	dfi.conf = conf
	return *dfi
}
