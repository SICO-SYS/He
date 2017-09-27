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

type AAAPrivateService struct{}

func Authentication(id string, signature string) (bool, int64) {
	r, err := mongo.FindOne(userDB, mongo.CollectionUserTokenName(), map[string]string{"id": id})
	if err != nil {
		return false, 201
	}
	prevTS, currentTS, nextTS := public.TimesPer30s()
	key, _ := r["key"].(string)
	if signature == public.EncryptWithSha256(key+prevTS) || signature == public.EncryptWithSha256(key+currentTS) || signature == public.EncryptWithSha256(key+nextTS) {
		return true, 0
	}
	return false, 0
}

func (a *AAAPrivateService) AuthenticationRPC(ctx context.Context, in *pb.AAATokenCall) (*pb.AAATokenBack, error) {
	isValid, errcode := Authentication(in.Id, in.Signature)
	if errcode != 0 {
		return &pb.AAATokenBack{Code: errcode}, nil
	}
	return &pb.AAATokenBack{Code: 0, IsValid: isValid}, nil
}

func (a *AAAPrivateService) AuthorizationRPC(ctx context.Context, in *pb.AAAServiceCall) (*pb.AAAServiceBack, error) {
	return &pb.AAAServiceBack{Code: 0}, nil
}

func (a *AAAPrivateService) AccountingRPC(ctx context.Context, in *pb.AAAEventCall) (*pb.AAAEventBack, error) {
	return &pb.AAAEventBack{Code: 0}, nil
}
