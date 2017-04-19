/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"fmt"
	"google.golang.org/grpc"
	"os"

	"github.com/SiCo-DevOps/Pb"
)

var (
	S   = grpc.NewServer()
	err error
)

func GenerateRand() string {
	data, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	defer data.Close()
	buf := make([]byte, 16)
	data.Read(buf)
	v := fmt.Sprintf("%X", buf)
	return v
}

func init() {
	pb.RegisterOpenServer(S, &Open{})
}
