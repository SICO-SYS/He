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

func Authentication(id string, signature string) bool {
	defer func() { recover() }()
	r := UserToken{}
	conn := userDB.Clone()
	defer conn.Close()
	conn.DB("SiCo").C("user.token").Find(mongo.Find("id", id)).One(&r)
	prevTS, currentTS, nextTS := public.TimesPer30s()
	if signature == public.EncryptWithSha256(r.Key+prevTS) || signature == public.EncryptWithSha256(r.Key+currentTS) || signature == public.EncryptWithSha256(r.Key+nextTS) {
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
