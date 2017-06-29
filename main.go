package main

import (
	"github.com/bububa/ProxyPool/crawler"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/storage"
	"github.com/bububa/ProxyPool/util"
	"log"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ipChan := make(chan models.IP, 2000)
	config := util.NewConfig()

	conn := storage.NewMysql(config)

	// Check the IPs in DB
	go func() {
		conn.CheckProxies()
	}()

	// Check the IPs in channel
	for i := 0; i < 50; i++ {
		go func() {
			for {
				conn.CheckProxy(<-ipChan)
			}
		}()
	}

	// Start getters to scraper IP and put it in channel
	for {
		total := conn.Count()
		log.Printf("Chan: %v, IP Count: %v\n", len(ipChan), total)
		if len(ipChan) < 100 {
			go run(ipChan)
		}
		time.Sleep(10 * time.Minute)
	}
}

func run(ipChan chan<- models.IP) {
	var wg sync.WaitGroup
	providers := []crawler.Provider{
		crawler.P5u,
		crawler.P66ip,
		crawler.PGoubanjia,
		crawler.PIp181,
		crawler.PKuaidaili,
		crawler.PProxylistplus,
		crawler.PXdaili,
		crawler.PXicidaili,
		crawler.PYoudaili,
	}
	for _, provider := range providers {
		wg.Add(1)
		go func(provider crawler.Provider) {
			ips, err := crawler.Fetch(provider)
			if err != nil {
				return
			}
			for _, ip := range ips {
				ipChan <- ip
			}
			wg.Done()
		}(provider)
	}
	wg.Wait()
	log.Println("All providers finished.")
}
