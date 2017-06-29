package crawler

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/request"
)

const (
	P5u_URL = "http://www.data5u.com/free/index.shtml"
)

type Crawler5u struct{}

// Data5u get ip from data5u.com
func (this Crawler5u) Fetch() (ips []models.IP, err error) {
	resp, _, errs := request.New().Get(P5u_URL).End()
	if errs != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	doc.Find("body > div.wlist > li:nth-child(2) > ul").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		addr := s.Find("ul:nth-child(" + node + ") > span:nth-child(1) > li").Text()
		port := s.Find("ul:nth-child(" + node + ") > span:nth-child(2) > li").Text()
		scheme := s.Find("ul:nth-child(" + node + ") > span:nth-child(4) > li").Text()
		ip := models.IP{
			Data: addr + ":" + port,
			Type: scheme,
		}
		ips = append(ips, ip)
	})
	return ips, nil
}
