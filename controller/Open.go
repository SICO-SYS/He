/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/mongo"
	. "github.com/SiCo-Ops/log"
	"github.com/SiCo-Ops/public"
)

var (
	id  string
	key string
)

type AAA_Open struct{}

func (o *AAA_Open) AAA_RegUser(ctx context.Context, in *pb.AAA_RegRequest) (*pb.AAA_APIKeypair, error) {
	email := in.Email
	notExist := false
	for i := 0; notExist == false; i++ {
		if i == 5 {
			notExist = true
			LogErrMsg(20, "controller.RegUser")
			return &pb.AAA_APIKeypair{Id: "", Key: ""}, nil
		}
		id = public.GenHexString()
		key = public.GenHexString()
		keypair := &UserKeypair{Id: id, Key: key, Email: email, Time: public.Now()}
		notExist = mongo.Mgo_Insert(mongo.MgoUserConn, keypair, "user.keypair")
		if notExist {
			auth := &UserAuth{Id: id}
			mongo.Mgo_Insert(mongo.MgoUserConn, auth, "user.auth")
		}
	}
	return &pb.AAA_APIKeypair{Id: id, Key: key}, nil
}
