/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
	. "github.com/SiCo-DevOps/log"
)

var (
	key    string
	secret string
)

type AAA_Open struct{}

func (o *AAA_Open) AAA_RegUser(ctx context.Context, in *pb.AAA_OpenRequest) (*pb.AAA_APIKeypair, error) {
	action := in.Request

	if action == "reg" {
		notExist := false
		for i := 0; notExist == false; i++ {
			if i == 5 {
				notExist = true
				LogErrMsg(20, "controller.RegUser")
				return &pb.AAA_APIKeypair{Key: "", Secret: ""}, nil
			}
			key = GenerateRand()
			secret = GenerateRand()
			enkey := Sha256Encrypt(key)
			keypair := &UserKeypair{Key: enkey, Secret: secret}
			auth := &UserAuth{Key: enkey}
			notExist = dao.Mgo_Insert(keypair, "user.keypair")
			dao.Mgo_Insert(auth, "user.auth")
		}
		return &pb.AAA_APIKeypair{Key: key, Secret: secret}, nil
	}
	return &pb.AAA_APIKeypair{}, nil
}
