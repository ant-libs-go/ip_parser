/* ######################################################################
# Author: (zhengfei@dianzhong.com)
# Created Time: 2021-08-09 17:09:33
# File Name: main_test.go
# Description:
####################################################################### */

package ip_parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ant-libs-go/util"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestDingDingText(t *testing.T) {
	dat, _ := ioutil.ReadFile("./qqwry.dat")

	util.ReadLine("./ip", func(line []byte) {
		fmt.Printf("%+v\n", NewIpParser(string(line), dat).Parse())
	})
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
