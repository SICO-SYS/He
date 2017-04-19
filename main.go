/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package main

import (
	"net"

	"github.com/SiCo-DevOps/He/controller"
)

func Run() {
	lis, _ := net.Listen("tcp", ":6666")
	controller.S.Serve(lis)
}

func main() {
	Run()
}
