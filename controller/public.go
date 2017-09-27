/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/getsentry/raven-go"
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
	if email == "" && phone == "" {
		return &pb.AAAGenerateTokenBack{Code: 1003}, nil
	}
	for i := 0; true; i++ {
		tokenID = public.GenerateHexString()
		tokenKey = public.GenerateHexString()
		token := &UserToken{tokenID, tokenKey, email, phone, public.Now()}
		err := mongo.Insert(userDB, mongo.CollectionUserTokenName(), token)
		if err != nil {
			if i >= 4 {
				raven.CaptureError(err, nil)
				return &pb.AAAGenerateTokenBack{Code: 201}, nil
			}
			continue
		}
		err = mongo.Insert(userDB, mongo.CollectionUserPolicyName(), &UserPolicy{Id: tokenID})
		if err != nil {
			raven.CaptureError(err, nil)
			return &pb.AAAGenerateTokenBack{Code: 201}, nil
		}
		break
	}
	return &pb.AAAGenerateTokenBack{Code: 0, Id: tokenID, Key: tokenKey}, nil
}
