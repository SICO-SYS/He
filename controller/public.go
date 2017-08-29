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
	"github.com/SiCo-Ops/public"
)

var (
	tokenID  string
	tokenKey string
)

type UserToken struct {
	Id    string "id"
	Key   string "key"
	Email string "email"
	Phone string "phone"
	Time  string "createtime"
}

type UserPolicy struct {
	Id          string   "id"
	OpsMember   []string "opsMember"
	AdminMember []string "adminMember"
	Group       []string "group"
}

type AAAPublicService struct{}

func (o *AAAPublicService) GenerateTokenRPC(ctx context.Context, in *pb.AAAGenerateTokenCall) (*pb.AAAGenerateTokenBack, error) {
	email := in.Email
	phone := in.Phone
	tokenExist := false
	for i := 0; tokenExist == false; i++ {
		if i == 5 {
			tokenExist = true
			return &pb.AAAGenerateTokenBack{Id: "", Key: ""}, nil
		}
		tokenID = public.GenerateHexString()
		tokenKey = public.GenerateHexString()
		token := &UserToken{tokenID, tokenKey, email, phone, public.Now()}
		tokenExist = mongo.Insert(mongo.UserConn, token, mongo.CollectionUserTokenName())
		if tokenExist == true {
			mongo.Insert(mongo.UserConn, &UserPolicy{Id: tokenID}, mongo.CollectionUserPolicyName())
		}
	}
	return &pb.AAAGenerateTokenBack{Id: tokenID, Key: tokenKey}, nil
}
