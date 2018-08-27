package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/adamluo159/mylog"
)

func main() {
	logpath := flag.String("log", "./log/proxyip.log", "log file")
	addr := flag.String("addr", ":6628", "proxyip server addr")
	flag.Parse()
	_, err := mylog.New(*logpath, mylog.LogDebug, -1, mylog.GB)
	if err != nil {
		panic(err)
	}

	mylog.Info("proxyip init .....")
	RequestProxyIps()
	go startProxyIpServer(*addr)
	mylog.Info("proxyip working .....")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	mylog.Info("proxyip shut down.....")
}
