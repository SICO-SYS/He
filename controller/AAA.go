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

type UserThirdparty struct {
	ID         string "id"
	CloudAlias string "name"
	CloudID    string "cloudid"
	CloudKey   string "cloudkey"
}

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

// func (s *AAA_Secret) AAA_ThirdKeypair(ctx context.Context, in *pb.AAA_ThirdpartyKey) (*pb.ResponseMsg, error) {
// 	if AAA(in.Apitoken.Id, in.Apitoken.Signature) {
// 		c := "user.cloud."
// 		switch strings.ToLower(in.Apitype) {
// 		case "aws":
// 			c += in.Apitype
// 		case "qcloud":
// 			c += in.Apitype
// 		case "aliyun":
// 			c += in.Apitype
// 		default:
// 			return &pb.ResponseMsg{Code: 2, Msg: "Not support cloud"}, nil
// 		}
// 		v := &UserThirdparty{in.Apitoken.Id, in.Name, in.Id, in.Key}
// 		ok := mongo.Mgo_Insert(mongo.MgoUserConn, v, c)
// 		if ok {
// 			return &pb.ResponseMsg{Code: 0}, nil
// 		}
// 		LogErrMsg(20, "controller.AAA_ThirdKeypair")
// 		return &pb.ResponseMsg{Code: 2, Msg: "Cannot Setup new keypair, maybe name exist"}, nil
// 	}
// 	return &pb.ResponseMsg{Code: 2, Msg: "AAA failed"}, nil
// }

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

// 	query := mongo.Mgo_Querys{"id": in.Apitoken.Id, "name": in.Name}
// 	result := query.Mgo_FindsOne(mongo.MgoUserConn, c)
// 	cloudid, ok := result["cloudid"].(string)
// 	cloudkey, _ := result["cloudkey"].(string)
// 	if !ok {
// 		return &pb.AAA_APIKeypair{Id: "getresult failed", Key: ""}, nil
// 	}
// 	return &pb.AAA_APIKeypair{Id: cloudid, Key: cloudkey}, nil
// }
