package crawler

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/phantomjs"
)

const (
	PKuaidaili_URL = "http://www.kuaidaili.com/proxylist/"
)

type CrawlerKuaidaili struct{}

func (this CrawlerKuaidaili) Fetch() (ips []models.IP, err error) {
	//create a fetcher which seems to a httpClient
	fetcher, err := phantomjs.NewFetcher(2016, nil)
	defer fetcher.ShutDownPhantomJSServer()
	if err != nil {
		return
	}
	//inject the javascript you want to run in the webpage just like in chrome console.
	jsScript := "function() {s=document.documentElement.outerHTML;document.write('<body></body>');document.body.innerText=s;}"
	//run the injected js_script at the end of loading html
	jsRunAt := phantomjs.RUN_AT_DOC_END
	//send httpGet request with injected js

	for i := 1; i <= 10; i++ {
		resp, err := fetcher.GetWithJS(PKuaidaili_URL+strconv.Itoa(i), jsScript, jsRunAt)
		if err != nil {
			return nil, err
		}

		//select search results by goquery
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.Content))
		if err != nil {
			return nil, err
		}
		doc.Find("#index_free_list > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
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
	}
	return ips, nil
}
