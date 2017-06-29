package storage

import (
	"errors"
	"github.com/bububa/ProxyPool/models"
	"github.com/bububa/ProxyPool/util"
	"github.com/bububa/mymysql/autorc"
	_ "github.com/bububa/mymysql/thrsafe" // Native engine
	"sync"
)

// Storage struct is used for storeing persistent data of alerts
type Mysql struct {
	*Storage
	table   string
	session *autorc.Conn
}

// NewStorage creates and returns new Storage instance
func NewMysql(conf *util.Config) *Mysql {
	mdb := autorc.New("tcp", "", conf.Mysql.Host, conf.Mysql.User, conf.Mysql.Passwd, conf.Mysql.DB)
	mdb.Register("set names utf8")
	return &Mysql{table: conf.Mysql.Table, session: mdb}
}

func (this *Mysql) Insert(ip models.IP) error {
	db := this.session
	_, _, err := db.Query(`INSERT INTO %s (ip, ip_type) VALUES ('%s', '%s')`, this.table, db.Escape(ip.Data), db.Escape(ip.Type))
	return err
}

func (this *Mysql) Get(value string) (models.IP, error) {
	db := this.session
	rows, _, err := db.Query(`SELECT ip, ip_type FROM %s WHERE ip='%s'`, this.table, db.Escape(value))
	if err != nil {
		return models.IP{}, err
	}
	if len(rows) == 0 {
		return models.IP{}, errors.New("not found")
	}
	return models.IP{Data: rows[0].Str(0), Type: rows[0].Str(1)}, nil
}

func (this *Mysql) Delete(ip models.IP) error {
	db := this.session
	_, _, err := db.Query(`DELETE FROM %s WHERE ip='%s'`, this.table, db.Escape(ip.Data))
	return err
}

func (this *Mysql) Update(ip models.IP) error {
	db := this.session
	_, _, err := db.Query(`INSERT INTO %s (ip, ip_type) VALUES ('%s', '%s') ON DUPLICATE KEY UPDATE ip_type=VALUES(ip_type)`, this.table, db.Escape(ip.Data), db.Escape(ip.Type))
	return err
}

func (this *Mysql) All() ([]models.IP, error) {
	db := this.session
	rows, _, err := db.Query(`SELECT ip, ip_type FROM %s`, this.table)
	if err != nil {
		return []models.IP{}, err
	}
	var ips []models.IP
	for _, row := range rows {
		ips = append(ips, models.IP{Data: row.Str(0), Type: row.Str(1)})
	}
	return ips, nil
}

func (this *Mysql) Random() (models.IP, error) {
	db := this.session
	rows, _, err := db.Query(`SELECT ip, ip_type FROM %s ORDER BY rand LIMIT 1`, this.table)
	if err != nil {
		return models.IP{}, err
	}
	if len(rows) == 0 {
		return models.IP{}, errors.New("not found")
	}
	return models.IP{Data: rows[0].Str(0), Type: rows[0].Str(1)}, nil
}

func (this *Mysql) Count() uint {
	db := this.session
	rows, _, err := db.Query(`SELECT COUNT(*) FROM %s`, this.table)
	if err != nil {
		return 0
	}
	return rows[0].Uint(0)
}

func (this *Mysql) CheckProxies() error {
	ips, err := this.All()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(v models.IP) {
			if !check(v) {
				this.Delete(v)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	return nil
}

func (this *Mysql) CheckProxy(ip models.IP) (bool, error) {
	if check(ip) {
		err := this.Insert(ip)
		return err == nil, err
	}
	return false, nil
}
