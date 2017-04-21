/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"crypto/sha256"
	"encoding/hex"
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

func Sha256Encrypt(v string) string {
	hash := sha256.New()
	hash.Write([]byte(v))
	return hex.EncodeToString(hash.Sum(nil))
}

func init() {
	pb.RegisterAAA_OpenServer(S, &AAA_Open{})
	pb.RegisterAAA_SecretServer(S, &AAA_Secret{})
}
