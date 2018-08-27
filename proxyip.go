package main

import (
	"sync"
)

var (
	httpips    []string
	httplocker sync.RWMutex

	httpsips    []string
	httpslocker sync.RWMutex

	ipcount int = 20
)

func delHttpIps(ips []string) {
	httplocker.Lock()
	defer httplocker.Unlock()

	for i := 0; i < len(ips); i++ {
		for j := 0; j < len(httpips); j++ {
			if httpips[j] == ips[i] {
				httpips = append(httpips[:j], httpips[j+1:]...)
				break
			}
		}
	}
}

func delHttpsIps(ips []string) {
	httpslocker.Lock()
	defer httpslocker.Unlock()

	for i := 0; i < len(ips); i++ {
		for j := 0; j < len(httpsips); j++ {
			if httpsips[j] == ips[i] {
				httpsips = append(httpsips[:j], httpsips[j+1:]...)
				break
			}
		}
	}
}

func addHttpIps(ips []string) {
	httplocker.Lock()
	defer httplocker.Unlock()
	httpips = append(httpips, ips...)
}

func addHttpsIps(ips []string) {
	httpslocker.Lock()
	defer httpslocker.Unlock()
	httpsips = append(httpsips, ips...)
}

func getHttpIps() []string {
	httplocker.RLock()
	defer httplocker.RUnlock()

	ips := make([]string, 0, ipcount)
	for i := 0; i < len(httpips); i++ {
		ips = append(ips, httpips[i])
	}
	return ips
}

func getHttpsIps() []string {
	httpslocker.RLock()
	defer httpslocker.RUnlock()

	ips := make([]string, 0, ipcount)
	for i := 0; i < len(httpsips); i++ {
		ips = append(ips, httpsips[i])
	}
	return ips
}
