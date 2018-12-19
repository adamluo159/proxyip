package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/adamluo159/mylog"
)

type IPFunc func(string)

func getDoc(req_url string) (*goquery.Document, error) {
	request, _ := http.NewRequest("GET", req_url, nil)
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36")
	client := &http.Client{
		Timeout: time.Duration(20 * time.Second),
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("response status:%d", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func XiCi(f func(string), wg *sync.WaitGroup) {
	xici_addr := "http://www.xicidaili.com/wt/"
	for i := 1; i <= 10; i++ {
		xici_url := fmt.Sprintf("%s%d", xici_addr, i)
		doc, err := getDoc(xici_url)
		if err != nil {
			mylog.Error("xici:%d request ip err:%v", xici_url, err)
			continue
		}
		doc.Find("#ip_list tbody .odd").Each(func(i int, context *goquery.Selection) {
			ip := context.Find("td").Eq(1).Text()
			port := context.Find("td").Eq(2).Text()
			proxy_ip := fmt.Sprintf("%s:%s", ip, port)
			f(proxy_ip)
			mylog.Debug("get xici ip:%s from page:%d", proxy_ip, i)
		})
	}
	wg.Done()
}

func KuaiDaiLi(f func(string), wg *sync.WaitGroup) {
	kuaidai_addr := "https://www.kuaidaili.com/free/inha/"
	for i := 1; i <= 10; i++ {
		kuaidai_url := fmt.Sprintf("%s%d/", kuaidai_addr, i)
		doc, err := getDoc(kuaidai_url)
		if err != nil {
			mylog.Error("kuaidai:%s request ip err:%v", kuaidai_url, err)
			continue
		}
		doc.Find("#list tbody tr").Each(func(i int, context *goquery.Selection) {
			ip := context.Find("td").Eq(0).Text()
			port := context.Find("td").Eq(1).Text()
			proxy_ip := fmt.Sprintf("%s:%s", ip, port)
			f(proxy_ip)
			mylog.Debug("get kuaidai ip:%s from page:%d", proxy_ip, i)
		})
	}
	wg.Done()
}

func Ip89(f func(string), wg *sync.WaitGroup) {
	ip89_addr := "http://www.89ip.cn"
	for i := 1; i <= 10; i++ {
		ip89_url := fmt.Sprintf("%s/index_%d.html", ip89_addr, i)
		doc, err := getDoc(ip89_url)
		if err != nil {
			mylog.Error("ip89_url:%s request ip err:%v", ip89_url, err)
			continue
		}
		doc.Find(".layui-table tbody tr").Each(func(i int, context *goquery.Selection) {
			ip := context.Find("td").Eq(0).Text()
			port := context.Find("td").Eq(1).Text()
			proxy_ip := fmt.Sprintf("%s:%s", ip, port)
			proxy_ip = strings.Replace(proxy_ip, "\t", "", -1)
			proxy_ip = strings.Replace(proxy_ip, "\n", "", -1)
			f(proxy_ip)
			mylog.Debug("get ip89 ip:%s from page:%d", proxy_ip, i)
		})
	}
	wg.Done()
}
