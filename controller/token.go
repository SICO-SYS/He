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

type AAATokenService struct{}

func Authentication(id string, signature string) (bool, int64) {
	r, err := mongo.FindOne(userDB, mongo.CollectionUserTokenName(), map[string]string{"id": id})
	if err != nil {
		return false, 201
	}
	if r == nil {
		return false, 1005
	}
	prevTS, currentTS, nextTS := public.TimesPer30s()
	key, _ := r["key"].(string)
	if signature == public.EncryptWithSha256(key+prevTS) || signature == public.EncryptWithSha256(key+currentTS) || signature == public.EncryptWithSha256(key+nextTS) {
		return true, 0
	}
	return false, 0
}

func (a *AAATokenService) GenerateRPC(ctx context.Context, in *pb.AAAGenerateTokenCall) (*pb.AAAGenerateTokenBack, error) {
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

func (a *AAATokenService) AuthenticationRPC(ctx context.Context, in *pb.AAATokenCall) (*pb.AAATokenBack, error) {
	isValid, errcode := Authentication(in.Id, in.Signature)
	if errcode != 0 {
		return &pb.AAATokenBack{Code: errcode}, nil
	}
	return &pb.AAATokenBack{Code: 0, IsValid: isValid}, nil
}

func (a *AAATokenService) AuthorizationRPC(ctx context.Context, in *pb.AAAServiceCall) (*pb.AAAServiceBack, error) {
	return &pb.AAAServiceBack{Code: 0}, nil
}

func (a *AAATokenService) AccountingRPC(ctx context.Context, in *pb.AAAEventCall) (*pb.AAAEventBack, error) {
	return &pb.AAAEventBack{Code: 0}, nil
}
