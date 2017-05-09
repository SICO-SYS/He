/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"
	"strings"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao/mongo"
	. "github.com/SiCo-DevOps/log"
	"github.com/SiCo-DevOps/public"
)

type UserKeypair struct {
	Id    string "id"
	Key   string "key"
	Email string "email"
	Time  string "createtime"
}

type UserAuth struct {
	Id          string   "id"
	QueryMember []string "query_member"
	AdminMember []string "admin_member"
	Group       []string "group"
}

type UserThirdparty struct {
	ID         string "id"
	CloudAlias string "name"
	CloudID    string "cloudid"
	CloudKey   string "cloudkey"
}

type AAA_Secret struct{}

func AAA(k string, t string) bool {
	defer func() { recover(); LogProduce("error", "Execute AAA failed") }()
	r := UserKeypair{}
	conn := mongo.MgoUserConn.Clone()
	defer conn.Close()
	conn.DB("SiCo").C("user.keypair").Find(mongo.Mgo_Find("id", k)).One(&r)
	sbefore, snow, safter := public.Per30sTimes()
	if t == public.Sha256Encrypt(r.Key+sbefore) || t == public.Sha256Encrypt(r.Key+snow) || t == public.Sha256Encrypt(r.Key+safter) {
		return true
	}
	LogErrMsg(50, "controller.AAA")
	return false
}

func (a *AAA_Secret) AAA_Auth(ctx context.Context, in *pb.AAA_APIToken) (*pb.ResponseMsg, error) {
	if AAA(in.Id, in.Signature) {
		return &pb.ResponseMsg{Code: 0}, nil
	}
	return &pb.ResponseMsg{Code: 2}, nil
}

func (s *AAA_Secret) AAA_ThirdKeypair(ctx context.Context, in *pb.AAA_ThirdpartyKey) (*pb.ResponseMsg, error) {
	if AAA(in.Apitoken.Id, in.Apitoken.Signature) {
		c := "user.cloud."
		switch strings.ToLower(in.Apitype) {
		case "aws":
			c += in.Apitype
		case "qcloud":
			c += in.Apitype
		case "aliyun":
			c += in.Apitype
		default:
			return &pb.ResponseMsg{Code: 2, Msg: "Not support cloud"}, nil
		}
		v := &UserThirdparty{in.Apitoken.Id, in.Name, in.Id, in.Key}
		ok := mongo.Mgo_Insert(mongo.MgoUserConn, v, c)
		if ok {
			return &pb.ResponseMsg{Code: 0}, nil
		}
		LogErrMsg(20, "controller.AAA_ThirdKeypair")
		return &pb.ResponseMsg{Code: 2, Msg: "Cannot Setup new keypair, maybe name exist"}, nil
	}
	return &pb.ResponseMsg{Code: 2, Msg: "AAA failed"}, nil
}

func (a *AAA_Secret) AAA_GetThirdKey(ctx context.Context, in *pb.AAA_ThirdpartyKey) (*pb.AAA_APIKeypair, error) {
	if !AAA(in.Apitoken.Id, in.Apitoken.Signature) {
		return &pb.AAA_APIKeypair{Id: "", Key: ""}, nil
	}
	c := "user.cloud."
	switch strings.ToLower(in.Apitype) {
	case "aws":
		c += in.Apitype
	case "qcloud":
		c += in.Apitype
	case "aliyun":
		c += in.Apitype
	default:
		return &pb.AAA_APIKeypair{Id: "cloudfailed", Key: ""}, nil
	}

	query := mongo.Mgo_Querys{"id": in.Apitoken.Id, "name": in.Name}
	result := query.Mgo_FindsOne(mongo.MgoUserConn, c)
	cloudid, ok := result["cloudid"].(string)
	cloudkey, _ := result["cloudkey"].(string)
	if !ok {
		return &pb.AAA_APIKeypair{Id: "getresult failed", Key: ""}, nil
	}
	return &pb.AAA_APIKeypair{Id: cloudid, Key: cloudkey}, nil
}
