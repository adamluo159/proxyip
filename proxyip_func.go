package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/adamluo159/mylog"
)

var (
	proxyfunc = []func(){
		Getxici,
	}
)

func RequestProxyIps() {
	for i := 0; i < len(proxyfunc); i++ {
		proxyfunc[i]()
	}
}

func getWebDoc(urls string, proxyUrl *url.URL) (*goquery.Document, error) {
	request, _ := http.NewRequest("GET", urls, nil)
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36")
	client := &http.Client{
		Timeout: time.Duration(20 * time.Second),
	}

	if proxyUrl != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("getWebDoc url:%s proxy:%+v err:%+v", urls, proxyUrl, err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("getWebDoc url:%s proxy:%+v statuscode:%+v", urls, proxyUrl, response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("getWebDoc url:%s proxy:%+v,  NewDocumentFromResponse err :%+v", urls, proxyUrl, err)
	}
	mylog.Debug("get web doc url:%s proxy:%+v", urls, proxyUrl)

	return doc, nil
}

func get(url_d string) (*goquery.Document, error) {
	var doc *goquery.Document = nil
	var err error
	var ips []string
	var del_ips []string = make([]string, 0, ipcount)

	isHttp := strings.Contains(url_d, "http")
	if !isHttp {
		isHttps := strings.Contains(url_d, "https")
		if !isHttps {
			return nil, fmt.Errorf("url err :%s", url_d)
		}
	}
	for {
		if isHttp {
			ips = getHttpIps()
		} else {
			ips = getHttpsIps()
		}

		if len(ips) == 0 {
			doc, err = getWebDoc(url_d, nil)
			if err != nil {
				mylog.Warn("%+v", err)
			}
			break
		} else {
			for i := 0; i < len(ips); i++ {
				purl, err := url.Parse(ips[i])
				if err != nil {
					mylog.Warn("parse url err:%+v proxy:%s", err, ips[i])
					del_ips = append(del_ips, ips[i])
					continue
				}
				doc, err = getWebDoc(url_d, purl)
				if err != nil {
					mylog.Warn("%+v", err)
					del_ips = append(del_ips, ips[i])
					continue
				}
				break
			}
		}

		if isHttp {
			delHttpIps(del_ips)
		} else {
			delHttpsIps(del_ips)
		}
	}

	return doc, err
}

func Getxici() {
	xici_addr := "http://www.xicidaili.com/wn/"
	https_ips := make([]string, 0)
	http_ips := make([]string, 0)

	for i := 1; i <= 20; i++ {
		xicipage := xici_addr + strconv.Itoa(i)
		doc, err := get(xicipage)
		if err != nil {
			mylog.Error("xici %+v", err)
			continue
		}
		doc.Find("#ip_list tbody .odd").Each(func(i int, context *goquery.Selection) {
			//地址
			ip := context.Find("td").Eq(1).Text()
			//端口
			port := context.Find("td").Eq(2).Text()
			//类型
			urlstr := context.Find("td").Eq(5).Text()

			addr := ip + ":" + port

			if urlstr == "HTTP" {
				addr = "http://" + addr
				http_ips = append(http_ips, addr)
			} else {
				addr = "https://" + addr
				https_ips = append(https_ips, addr)
			}
			mylog.Debug("xici get proxy index:%d url:%+v", i, addr)
		})
		addHttpIps(http_ips)
		addHttpsIps(https_ips)
		http_ips = http_ips[:0]
		https_ips = https_ips[:0]
	}
}
