package main

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/adamluo159/mylog"
)

type ProxyIp struct {
	Addr   string `json:"addr"`
	Http   bool   `json:"http"`
	usable bool
}

var (
	proxy_map    map[string]*ProxyIp = make(map[string]*ProxyIp)
	proxy_locker sync.RWMutex
	ipcount      int = 20
	cancel_func  context.CancelFunc
	close        chan bool = make(chan bool)
	fileName     string    = "proxyfile"
)

func Run(serverAddr string) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				saveIpsToFile()
				close <- true
				return
			case <-time.After(time.Minute * 5):
				mylog.Info("begin check proxy valid.. len:%d", len(proxy_map))
				tmp_proxy_map := make(map[string]*ProxyIp)
				proxy_locker.Lock()
				for _, v := range proxy_map {
					tmp_proxy_map[v.Addr] = v
				}
				proxy_locker.Unlock()

				for _, tmp_v := range tmp_proxy_map {
					tmp_v.usable = checkProxy("http://" + tmp_v.Addr)
					if !tmp_v.usable {
						proxy_locker.Lock()
						delete(proxy_map, tmp_v.Addr)
						proxy_locker.Unlock()
					}
				}
				mylog.Info("end check proxy valid.. len:%d", len(proxy_map))
				if len(proxy_map) < 5 {
					RequestProxyIps()
				}
			}
		}
	}(ctx)

	readIpsFromFile()
	if len(proxy_map) == 0 {
		RequestProxyIps()
	}

	go startProxyIpServer(serverAddr)
	cancel_func = cancel
}

func Destory() {
	mylog.Info("proxyip begin destory.")
	cancel_func()
	c := <-close
	proxy_srv.Shutdown(nil)
	mylog.Info("proxyip done ", c)
}

func delIps(ips []string) {
	proxy_locker.Lock()
	defer proxy_locker.Unlock()
	for i := 0; i < len(ips); i++ {
		ip := ips[i]
		_, ok := proxy_map[ip]
		if ok {
			proxy_map[ip].usable = false
		}
	}
}

func addIp(addr string, isHttp bool) {
	proxy_locker.Lock()
	defer proxy_locker.Unlock()
	proxy_map[addr] = &ProxyIp{
		Addr:   addr,
		Http:   isHttp,
		usable: true,
	}
}
func addProxyIp(proxy *ProxyIp) {
	proxy_locker.Lock()
	defer proxy_locker.Unlock()
	proxy.usable = true
	proxy_map[proxy.Addr] = proxy
}

func getIps() []*ProxyIp {
	proxy_locker.RLock()
	defer proxy_locker.RUnlock()

	i := 0
	proxys := make([]*ProxyIp, 0, ipcount)
	for _, v := range proxy_map {
		if v.usable {
			proxys = append(proxys, v)
			if i > ipcount {
				break
			}
			i++
		}
	}
	return proxys
}

func saveIpsToFile() {
	mylog.Info("begin saveIpsToFile.")
	proxy_locker.Lock()
	defer proxy_locker.Unlock()
	f, err := os.Create(fileName)
	if err != nil {
		mylog.Warn("saveIpsToFile %+v", err)
	}
	defer f.Close()

	buf := bufio.NewWriter(f)
	for _, v := range proxy_map {
		b, err := json.Marshal(v)
		if err != nil {
			mylog.Warn("saveIpsToFile err:%+v proxy:%+v", err, v)
			continue
		}
		b = append(b, byte('\n'))
		_, err = buf.Write(b)
		if err != nil {
			mylog.Warn("saveIpsToFile write into file err:%+v proxy:%+v", err, v)
		}
		mylog.Debug("err:%+v save %+v", err, string(b))
	}
	err = buf.Flush()
	if err != nil {
		mylog.Info("saveIpsToFile fail. ")
	}
	mylog.Info("end saveIpsToFile. ")
}

func readIpsFromFile() {
	mylog.Info("begin readIpsFromFile.")

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		mylog.Warn("readIpsFromFile read %s err:%+v", fileName, err)
		return
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			break
		}
		proxy := &ProxyIp{}
		err = json.Unmarshal(line, proxy)
		if err != nil {
			mylog.Warn("readIpsFromFile json unmarshal err:%+v line:%s", err, line)
			continue
		}
		mylog.Debug("read %+v", proxy)
		addProxyIp(proxy)
	}

	mylog.Info("end readIpsFromFile.len:%d", len(proxy_map))

}
