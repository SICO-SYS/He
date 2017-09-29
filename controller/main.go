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
	"github.com/SiCo-Ops/cfg"
	"github.com/SiCo-Ops/dao/mongo"
)

const (
	configPath string = "config.json"
)

var (
	config            cfg.ConfigItems
	RPCServer         = grpc.NewServer()
	userDB, userDBErr = mongo.NewDial()
)

func ServePort() string {
	return config.RpcHePort
}

func init() {
	data, err := cfg.ReadFilePath(configPath)
	if err != nil {
		data = cfg.ReadConfigServer()
		if data == nil {
			log.Fatalln("config.json not exist and configserver was down")
		}
	}
	cfg.Unmarshal(data, &config)

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
