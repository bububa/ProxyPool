package crawler

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/request"

	"strings"
)

const (
	PIp181_URL = "http://www.ip181.com"
)

type CrawlerIp181 struct{}

// Data5u get ip from data5u.com
func (this CrawlerIp181) Fetch() (ips []models.IP, err error) {
	resp, _, errs := request.New().Get(PIp181_URL).End()
	if errs != nil {
		return nil, errs[0]
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	doc.Find("body > div:nth-child(3) > div.panel.panel-info > div.panel-body > div > div:nth-child(2) > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		sf, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(1)").Html()
		ff, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(2)").Html()
		hh, _ := s.Find("tr:nth-child(" + node + ") > td:nth-child(4)").Html()
		ip := models.IP{
			Data: sf + ":" + ff,
			Type: strings.ToLower(hh),
		}
		ips = append(ips, ip)
	})
	return ips[1:], nil
}
