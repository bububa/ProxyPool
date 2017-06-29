package crawler

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/request"
)

const (
	PProxylistplus_URL = "https://list.proxylistplus.com/"
)

type CrawlerProxylistplus struct{}

func (this CrawlerProxylistplus) Fetch() (ips []models.IP, err error) {
	_, body, errs := request.New().Get(PProxylistplus_URL).End()
	if errs != nil {
		return nil, errs[0]
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	doc.Find("#page > table.bg > tbody > tr").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		addr, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(2)").Html()
		port, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(3)").Html()
		scheme, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(7)").Html()
		if scheme == "yes" {
			scheme = "http,https"
		} else if scheme == "no" {
			scheme = "http"
		}
		ip := models.IP{
			Data: addr + ":" + port,
			Type: scheme,
		}
		ips = append(ips, ip)
	})
	return ips[2:], nil
}
