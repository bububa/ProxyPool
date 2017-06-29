package storage

import (
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/request"
)

const (
	HEALTH_CHECK_URL = "http://httpbin.org/get"
)

type Storage struct{}

// func (this *Storage) Random() (ip models.IP, err error) {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	ips, err := this.All()
// 	if err != nil {
// 		return ip, err
// 	}
// 	total := len(ips)

// 	return ips[r.Intn(total)]
// }

// CheckProxies to check the ip in DB

// CheckIP is to check the ip work or not
func check(ip models.IP) bool {
	resp, _, err := request.New().Proxy("http://" + ip.Data).Get(HEALTH_CHECK_URL).End()
	if err != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}
