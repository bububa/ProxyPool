package crawler

import (
	"github.com/bububa/ProxyPool/models"
)

type Provider uint

const (
	P5u            Provider = 1
	P66ip          Provider = 2
	PGoubanjia     Provider = 3
	PIp181         Provider = 4
	PKuaidaili     Provider = 5
	PProxylistplus Provider = 6
	PXdaili        Provider = 7
	PXicidaili     Provider = 8
	PYoudaili      Provider = 9
)

type Crawler interface {
	Fetch() ([]*models.IP, error)
}

func Fetch(provider Provider) ([]models.IP, error) {
	switch provider {
	case P5u:
		return Crawler5u{}.Fetch()
	case P66ip:
		return Crawler66ip{}.Fetch()
	case PGoubanjia:
		return CrawlerGoubanjia{}.Fetch()
	case PIp181:
		return CrawlerIp181{}.Fetch()
	case PKuaidaili:
		return CrawlerKuaidaili{}.Fetch()
	case PProxylistplus:
		return CrawlerProxylistplus{}.Fetch()
	case PXdaili:
		return CrawlerXdaili{}.Fetch()
	case PXicidaili:
		return CrawlerXicidaili{}.Fetch()
	case PYoudaili:
		return CrawlerYoudaili{}.Fetch()
	}
	return nil, nil
}
