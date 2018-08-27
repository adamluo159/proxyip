package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adamluo159/mylog"
)

func startProxyIpServer(addr string) {
	http.HandleFunc("/getips", GetIpHandler)
	http.HandleFunc("/delips", DelIpHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)

	}

}

type ProxyIp struct {
	HttpIps  []string `json:"http_ips"`
	HttpsIps []string `json:"https_ips"`
}

func GetIpHandler(w http.ResponseWriter, r *http.Request) {
	p := &ProxyIp{
		HttpsIps: getHttpsIps(),
		HttpIps:  getHttpIps(),
	}
	defer r.Body.Close()

	b, _ := json.Marshal(p)
	fmt.Fprintf(w, string(b))

}

func DelIpHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
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
}
