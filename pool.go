package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/adamluo159/mylog"
)

const (
	ip_update_Time time.Duration = time.Minute * 2
	max_ips        int           = 200
)

type pool struct {
	ip_lst     []string
	grap_funcs []func(func(string), *sync.WaitGroup)
	sync.Mutex
}

func NewPool() *pool {
	p := &pool{
		ip_lst:     make([]string, 0, max_ips),
		grap_funcs: []func(func(string), *sync.WaitGroup){XiCi, KuaiDaiLi},
	}

	go p.webServer(":3301")
	go p.updateIps()
	return p
}

func (p *pool) webServer(addr string) {
	http.HandleFunc("/getips", p.GetIpHandler)
	http.HandleFunc("/delips", p.DelIpHandler)
	if err := http.ListenAndServe(addr, nil); err == nil {
		mylog.Info("%s http server start success ", addr)
	} else {
		panic(err)
	}
}

func (p *pool) updateIps() {
	wg := &sync.WaitGroup{}
	for {
		iplen := len(p.ip_lst)
		if iplen < max_ips {
			mylog.Info("time to grap new  proxyips..")
			wg.Add(len(p.grap_funcs))
			for i := 0; i < len(p.grap_funcs); i++ {
				go p.grap_funcs[i](p.AddIp, wg)
			}
			wg.Wait()
		}
		time.Sleep(ip_update_Time)
	}
}

func (p *pool) AddIp(ip string) {
	p.Lock()
	defer p.Unlock()
	p.ip_lst = append(p.ip_lst, ip)
}

func (p *pool) DelIpHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return
	}

	dels := make([]string, 0)
	err = json.Unmarshal(b, &dels)
	if err != nil {
		return
	}

	p.Lock()
	for j := 0; j < len(dels); j++ {
		for i := 0; i < len(p.ip_lst); i++ {
			v := p.ip_lst[i]
			if v == dels[j] {
				p.ip_lst = append(p.ip_lst[:i], p.ip_lst[i+1:]...)
				break
			}
		}
	}
	p.Unlock()

	fmt.Fprintf(w, string("ok"))
}

func (p *pool) GetIpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	p.Lock()
	defer p.Unlock()

	iplen := len(p.ip_lst)
	if iplen == 0 {
		return
	}

	ips := make([]string, 0, iplen)
	for i := 0; i < iplen; i++ {
		n := rand.Intn(iplen)
		ips = append(ips, p.ip_lst[n])
	}

	b, _ := json.Marshal(ips)
	fmt.Fprintf(w, string(b))
}

func main() {
	mylog.New("./log/proxyip.log", mylog.LogDebug, time.Hour*24, mylog.GB)
	NewPool()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
