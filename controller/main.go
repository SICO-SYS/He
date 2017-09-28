/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/getsentry/raven-go"
	"google.golang.org/grpc"
	"log"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/cfg/v2"
	"github.com/SiCo-Ops/dao/mongo"
	"github.com/SiCo-Ops/dao/redis"
)

var (
	config            cfg.ConfigItems
	configPool        = redis.NewPool()
	RPCServer         = grpc.NewServer()
	userDB, userDBErr = mongo.NewDial()
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

	configPool = redis.InitPool(config.RedisConfigHost, config.RedisConfigPort, config.RedisConfigAuth)
	configs, err := redis.Hgetall(configPool, "system.config")
	if err != nil {
		log.Fatalln(err)
	}
	cfg.Map2struct(configs, &config)

	userDB, userDBErr = mongo.InitDial(config.MongoUserAddress, config.MongoUserUsername, config.MongoUserPassword)
	if userDBErr != nil {
		log.Fatalln(userDBErr)
	}
	err = mongo.AAAEnsureIndexes(userDB)
	if err != nil {
		log.Fatalln(err)
	}

	pb.RegisterAAATokenServiceServer(RPCServer, &AAATokenService{})

	if config.SentryHeStatus == "active" && config.SentryHeDSN != "" {
		raven.SetDSN(config.SentryHeDSN)
	}
}
