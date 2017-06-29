package crawler

import (
	"strings"

	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/request"
)

const (
	P66ip_URL = "http://www.66ip.cn/mo.php?tqsl=100&submit=%CC%E1++%C8%A1"
)

type Crawler66ip struct{}

// IP66 get ip from 66ip.cn
func (this Crawler66ip) Fetch() (ips []models.IP, err error) {
	_, body, errs := request.New().Get(P66ip_URL).End()
	if errs != nil {
		return nil, errs[0]
	}

	body = strings.Split(body, "c.js'></script>")[1]
	body = strings.Split(body, "</div>")[0]
	body = strings.TrimSpace(body)
	body = strings.Replace(body, "	", "", -1)
	temp := strings.Split(body, "<br />")
	for index := 0; index < len(temp[:len(temp)-1]); index++ {
		ip := models.IP{
			Data: strings.TrimSpace(temp[index]),
			Type: "http",
		}
		ips = append(ips, ip)
	}
	return ips, nil
}
