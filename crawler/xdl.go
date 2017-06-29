package crawler

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/phantomjs"
)

const (
	PXdaili_URL = "http://www.xdaili.cn/freeproxy.html"
)

type CrawlerXdaili struct{}

func (this CrawlerXdaili) Fetch() (ips []models.IP, err error) {
	fetcher, err := phantomjs.NewFetcher(2015, nil)
	defer fetcher.ShutDownPhantomJSServer()
	if err != nil {
		return nil, err
	}
	jsScript := "function() {s=document.documentElement.outerHTML;document.write('<body></body>');document.body.innerText=s;}"
	jsRunAt := phantomjs.RUN_AT_DOC_END
	resp, err := fetcher.GetWithJS(PXdaili_URL, jsScript, jsRunAt)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(strings.Replace(strings.Replace(resp.Content, "&lt;", "<", -1), "&gt;", ">", -1)))
	if err != nil {
		return nil, err
	}
	doc.Find("#target > tr").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		addr, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(1)").Html()
		port, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(2)").Html()
		scheme, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4)").Html()
		scheme = strings.Replace(strings.ToLower(scheme), "/", ",", -1)
		ip := models.IP{
			Data: addr + ":" + port,
			Type: scheme,
		}
		ips = append(ips, ip)
	})
	return ips, nil
}
