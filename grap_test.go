package main

import (
	"sync"
	"testing"
	"time"

	"github.com/adamluo159/mylog"
)

func Test_Xici(t *testing.T) {
	mylog.New("./log/test.log", mylog.LogDebug, time.Hour, mylog.GB)

	f := func(str string) {
	}
	wg := &sync.WaitGroup{}
	wg.Add(3)
	//XiCi(f, wg)
	//KuaiDaiLi(f, wg)
	Ip89(f, wg)
}
