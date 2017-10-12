/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"testing"
)

func Test_Authentication(t *testing.T) {
	_, errcode := Authentication("test", "test")
	if errcode != 1005 {
		t.Error(errcode)
	}
}

func Benchmark_Authentication(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, errcode := Authentication("test", "signature")
		if errcode != 1005 {
			b.Error(errcode)
		}
	}
}
