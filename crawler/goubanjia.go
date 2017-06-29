package crawler

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/request"
	"regexp"
	"strings"
)

const (
	PGoubanjia_URL = "http://www.goubanjia.com/free/gngn/index"
)

type CrawlerGoubanjia struct{}

func (this CrawlerGoubanjia) Fetch() (ips []models.IP, err error) {
	for i := 1; i <= 10; i++ {
		resp, _, errs := request.New().Get(PGoubanjia_URL + strconv.Itoa(i) + ".shtml").End()
		if errs != nil {
			return nil, errs[0]
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		doc.Find("#list > table > tbody > tr").Each(func(_ int, s *goquery.Selection) {
			sf, _ := s.Find(".ip").Html()
			tee := regexp.MustCompile("<pstyle=\"display:none;\">.?.?</p>").ReplaceAllString(strings.Replace(sf, " ", "", -1), "")
			re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
			ip := models.IP{
				Data: re.ReplaceAllString(tee, ""),
				Type: s.Find("td:nth-child(3) > a").Text(),
			}
			ips = append(ips, ip)
		})
	}
	return ips, nil
}
