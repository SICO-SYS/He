/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"google.golang.org/grpc"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
)

var (
	S   = grpc.NewServer()
	err error
)

func init() {
	defer func() {
		recover()
	}()
	dao.AAA_ensureIndexes()
	pb.RegisterAAA_OpenServer(S, &AAA_Open{})
	pb.RegisterAAA_SecretServer(S, &AAA_Secret{})
}
