/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"
	// "math"
	// "strconv"
	"strings"
	// "time"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
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
	conn := dao.MgoUserConn.Clone()
	defer conn.Close()
	conn.DB("SiCo").C("user.keypair").Find(dao.Mgo_Find("id", k)).One(&r)
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
		ok := dao.Mgo_Insert(v, c)
		if ok {
			return &pb.ResponseMsg{Code: 0}, nil
		}
		LogErrMsg(20, "controller.AAA_ThirdKeypair")
		return &pb.ResponseMsg{Code: 2, Msg: "Cannot Setup new keypair, maybe name exist"}, nil
	}
	return &pb.ResponseMsg{Code: 2, Msg: "AAA failed"}, nil
}
