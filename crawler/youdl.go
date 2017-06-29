package crawler

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/request"
)

const (
	PYoudaili_URL = "http://www.youdaili.net/Daili/http/"
)

type CrawlerYoudaili struct{}

func (this CrawlerYoudaili) Fetch() (ips []models.IP, err error) {
	_, body, errs := request.New().Get(PYoudaili_URL).End()
	if errs != nil {
		return nil, errs[0]
	}
	do, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	URL, _ := do.Find("body > div.con.PT20 > div.conl > div.lbtc.l > div.chunlist > ul > li:nth-child(1) > p > a").Attr("href")
	_, content, errs := request.New().Get(URL).End()
	if errs != nil {
		return nil, errs[0]
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	doc.Find(".content p").Each(func(_ int, s *goquery.Selection) {
		c := strings.Split(s.Text(), "@")
		ip := models.IP{
			Data: c[0],
			Type: strings.ToLower(strings.Split(c[1], "#")[0]),
		}
		ips = append(ips, ip)
	})
	return ips, nil
}
