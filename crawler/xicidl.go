package crawler

import (
	"regexp"
	"strings"

	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/phantomjs"
)

const (
	PXicidaili_URL = "http://www.xicidaili.com/nn/"
)

type CrawlerXicidaili struct{}

func (this CrawlerXicidaili) Fetch() (ips []models.IP, err error) {

	fetcher, err := phantomjs.NewFetcher(2017, nil)
	defer fetcher.ShutDownPhantomJSServer()
	if err != nil {
		return nil, err
	}
	jsScript := "function() {s=document.documentElement.outerHTML;document.write('<body></body>');document.body.innerText=s;}"
	jsRunAt := phantomjs.RUN_AT_DOC_END
	resp, err := fetcher.GetWithJS(PXicidaili_URL, jsScript, jsRunAt)
	if err != nil {
		return nil, err
	}
	re, _ := regexp.Compile("<td>(\\d+\\.){3}\\d+</td>.+?(\\d{2,4})</td>")
	temp := re.FindAllString(strings.Replace(strings.Replace(resp.Content, "&lt;", "<", -1), "&gt;", ">", -1), -1)

	for _, v := range temp {
		v = strings.Replace(v, "<td>", "", -1)
		v = strings.Replace(v, "</td>", "", -1)
		v = strings.Replace(v, " ", "", -1)
		v = strings.Replace(v, "<br>", ":", -1)
		ip := models.IP{
			Data: v,
			Type: "http",
		}
		ips = append(ips, ip)
	}
	return ips, nil
}
