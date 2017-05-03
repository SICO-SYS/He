/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package main

import (
	"net"

	"github.com/SiCo-DevOps/He/controller"
	"github.com/SiCo-DevOps/cfg"
)

func Run() {
	lis, _ := net.Listen("tcp", cfg.Config.Rpc.He)
	controller.S.Serve(lis)
}

func main() {
	Run()
}
