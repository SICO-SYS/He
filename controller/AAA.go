/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"
	// "strings"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/mongo"
	"github.com/SiCo-Ops/public"
)

type AAAPrivateService struct{}

func Authentication(k string, t string) bool {
	defer func() { recover() }()
	r := UserToken{}
	conn := mongo.MgoUserConn.Clone()
	defer conn.Close()
	conn.DB("SiCo").C("user.token").Find(mongo.MgoFind("id", k)).One(&r)
	prevTS, currentTS, nextTS := public.TimesPer30s()
	if t == public.EncryptWithSha256(r.Key+prevTS) || t == public.EncryptWithSha256(r.Key+currentTS) || t == public.EncryptWithSha256(r.Key+nextTS) {
		return true
	}
	return false
}

func (a *AAAPrivateService) AuthenticationRPC(ctx context.Context, in *pb.AAATokenCall) (*pb.AAATokenBack, error) {
	if Authentication(in.Id, in.Signature) {
		return &pb.AAATokenBack{Valid: true}, nil
	}
	return &pb.AAATokenBack{Valid: false}, nil
}

// func (a *AAA_Secret) AAA_GetThirdKey(ctx context.Context, in *pb.AAA_ThirdpartyKey) (*pb.AAA_APIKeypair, error) {
// 	if !AAA(in.Apitoken.Id, in.Apitoken.Signature) {
// 		return &pb.AAA_APIKeypair{Id: "", Key: ""}, nil
// 	}
// 	c := "user.cloud."
// 	switch strings.ToLower(in.Apitype) {
// 	case "aws":
// 		c += in.Apitype
// 	case "qcloud":
// 		c += in.Apitype
// 	case "aliyun":
// 		c += in.Apitype
// 	default:
// 		return &pb.AAA_APIKeypair{Id: "cloudfailed", Key: ""}, nil
// 	}

// }
