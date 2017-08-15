/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/getsentry/raven-go"
	"google.golang.org/grpc"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/cfg"
	"github.com/SiCo-Ops/dao/mongo"
)

var (
	config    = cfg.Config
	RPCServer = grpc.NewServer()
)

func init() {
	defer func() {
		recover()
	}()

	mongo.AAAEnsureIndexes()
	pb.RegisterAAAPublicServiceServer(RPCServer, &AAAPublicService{})
	pb.RegisterAAAPrivateServiceServer(RPCServer, &AAAPrivateService{})

	if config.Sentry.Enable {
		raven.SetDSN(config.Sentry.DSN)
	}
}
