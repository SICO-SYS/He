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
	"github.com/SiCo-Ops/cfg/v2"
	"github.com/SiCo-Ops/dao/mongo"
	"github.com/SiCo-Ops/dao/redis"
)

var (
	config            cfg.ConfigItems
	configPool        = redis.Pool("", "", "")
	RPCServer         = grpc.NewServer()
	userDB, userDBErr = mongo.Dial("", "", "")
)

func ServePort() string {
	return config.RpcHePort
}

func init() {
	defer func() {
		recover()
	}()
	data := cfg.ReadLocalFile()

	if data != nil {
		cfg.Unmarshal(data, &config)
	}

	configPool = redis.Pool(config.RedisConfigHost, config.RedisConfigPort, config.RedisConfigAuth)
	configs, _ := redis.Hgetall(configPool, "system.config")
	cfg.Map2struct(configs, &config)

	userDB, userDBErr = mongo.Dial(config.MongoUserAddress, config.MongoUserUsername, config.MongoUserPassword)
	mongo.AAAEnsureIndexes(userDB)
	pb.RegisterAAAPublicServiceServer(RPCServer, &AAAPublicService{})
	pb.RegisterAAAPrivateServiceServer(RPCServer, &AAAPrivateService{})

	if config.SentryHeStatus == "active" && config.SentryHeDSN != "" {
		raven.SetDSN(config.SentryHeDSN)
	}
}
