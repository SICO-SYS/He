/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"golang.org/x/net/context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
	. "github.com/SiCo-DevOps/log"
)

type UserKeypair struct {
	Key    string "key"
	Secret string "secret"
}

type UserAuth struct {
	Key         string   "key"
	QueryMember []string "query_member"
	AdminMember []string "admin_member"
	Group       []string "group"
}

type UserThirdparty struct {
	Key   string "key"
	ID    string "id"
	Token string "token"
}

type AAA_Secret struct{}

func AAA(k string, t string) bool {
	defer func() { recover(); LogProduce("error", "Execute AAA failed") }()
	r := UserKeypair{}
	conn := dao.MgoUserConn.Clone()
	defer conn.Close()
	conn.DB("SiCo").C("user.keypair").Find(dao.Mgo_Find("key", k)).One(&r)
	ts := time.Now().Unix()
	now := float64(ts / 30)
	stnow := strconv.Itoa(int(math.Floor(now)))
	stnowf := strconv.Itoa(int(math.Floor(now + 1)))
	stnows := strconv.Itoa(int(math.Floor(now - 1)))
	if t == Sha256Encrypt(r.Secret+stnow) || t == Sha256Encrypt(r.Secret+stnows) || t == Sha256Encrypt(r.Secret+stnowf) {
		return true
	}
	LogErrMsg(50, "controller.AAA")
	return false
}

func (a *AAA_Secret) AAA_Auth(ctx context.Context, in *pb.AAA_APIToken) (*pb.ResponseMsg, error) {
	if AAA(in.Key, in.Token) {
		return &pb.ResponseMsg{Code: 0}, nil
	}
	return &pb.ResponseMsg{Code: 2}, nil
}

func (s *AAA_Secret) AAA_ThirdKeypair(ctx context.Context, in *pb.AAA_ThirdpartyKey) (*pb.ResponseMsg, error) {
	if AAA(in.Apitoken.Key, in.Apitoken.Token) {
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
		v := &UserThirdparty{in.Apitoken.Key, in.Id, in.Key}
		ok := dao.Mgo_Insert(v, c)
		if ok {
			return &pb.ResponseMsg{Code: 0}, nil
		}
		LogErrMsg(20, "controller.AAA_ThirdKeypair")
	}
	return &pb.ResponseMsg{Code: 2, Msg: "AAA failed"}, nil
}
