package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

func getWebDoc(urls, proxyUrl string) (*goquery.Document, error) {
	request, _ := http.NewRequest("GET", urls, nil)
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36")
	client := &http.Client{
		Timeout: time.Duration(20 * time.Second),
	}

	if proxyUrl != "" {
		purl, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(purl),
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
	var del_ips []string = make([]string, 0, ipcount)

	//	isHttp := strings.Contains(url_d, "http")
	//	if !isHttp {
	//		isHttps := strings.Contains(url_d, "https")
	//		if !isHttps {
	//			return nil, fmt.Errorf("url err :%s", url_d)
	//		}
	//	}
	for {
		ips := getIps()
		if len(ips) == 0 {
			doc, err = getWebDoc(url_d, "")
			if err != nil {
				mylog.Warn("%+v", err)
			}
			break
		} else {
			for i := 0; i < len(ips); i++ {
				doc, err = getWebDoc(url_d, ips[i].Addr)
				if err != nil {
					mylog.Warn("%+v", err)
					del_ips = append(del_ips, ips[i].Addr)
					continue
				}
				break
			}
		}
		delIps(del_ips)
		del_ips = del_ips[:0]
	}

	return doc, err
}

func Getxici() {
	xici_addr := "http://www.xicidaili.com/wn/"
	for i := 1; i <= 20; i++ {
		xicipage := xici_addr + strconv.Itoa(i)
		doc, err := get(xicipage)
		if err != nil {
			mylog.Error("xici %+v", err)
			continue
		}
		doc.Find("#ip_list tbody .odd").Each(func(i int, context *goquery.Selection) {
			ip := context.Find("td").Eq(1).Text()
			port := context.Find("td").Eq(2).Text()
			urlstr := context.Find("td").Eq(5).Text()
			addr := ip + ":" + port

			addIp(addr, urlstr == "HTTP")
			mylog.Debug("xici get proxy index:%d type:%s url:%+v", i, urlstr, addr)
		})
	}
}

func checkProxy(proxyurl string) bool {
	checkurl := "http://www.baidu.com"
	_, err := getWebDoc(checkurl, proxyurl)
	if err != nil {
		return false
	}
	return true
}
