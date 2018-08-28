package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/adamluo159/mylog"
)

var proxy_srv *http.Server = &http.Server{
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}

func startProxyIpServer(addr string) {
	proxy_srv.Addr = addr
	http.HandleFunc("/getips", GetIpHandler)
	http.HandleFunc("/delips", DelIpHandler)
	mylog.Info("startProxyIpServer addr:%s", addr)
	proxy_srv.ListenAndServe()
}

func GetIpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ips := getIps()
	b, _ := json.Marshal(ips)
	fmt.Fprintf(w, string(b))
}

func DelIpHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		mylog.Warn("read delip http err:%+v, body:%s", err, string(b))
		return
	}
	p := &ProxyIp{}
	err = json.Unmarshal(b, p)
	if err != nil {
		mylog.Warn("read delip http jsonerr:%+v, body:%s", err, string(b))
		return
	}
	fmt.Fprintf(w, string("ok"))
}
