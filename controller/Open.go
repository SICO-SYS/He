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

type Open struct{}

func (o *Open) RegUser(ctx context.Context, in *pb.OpenRequest) (*pb.APIKeypair, error) {
	action := in.Request

	if action == "reg" {
		notExist := false
		for i := 0; notExist == false; i++ {
			if i == 5 {
				notExist = true
				LogErrMsg(20, "controller.RegUser")
				return &pb.APIKeypair{Key: "", Secret: ""}, nil
			}
			key = GenerateRand()
			secret = GenerateRand()
			data := &dao.UserKeypair{Key: key, Secret: secret}
			notExist = data.Insert()
		}
		return &pb.APIKeypair{Key: key, Secret: secret}, nil
	}
	return &pb.APIKeypair{}, nil
}
