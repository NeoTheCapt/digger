package utils

import (
	"database/sql"
	"digger/common"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/pool.v3"
	"log"
	"net"
	"strings"
	"time"
)

type mySQLScan struct {
	ips       []net.IP
	port      int16
	creds     []credential
	threadNum uint
	timeout   int
}

func InitMySQLScan(ips []net.IP, port int16, creds []string, threadsNum uint, timeout int) mySQLScan {
	real_creds := []credential{}

	for _, i := range creds {

		username := strings.Split(i, ":")[0]
		password := strings.Split(i, ":")[1]
		real_creds = append(real_creds, InitCred(username, password))
	}
	return mySQLScan{
		ips:       ips,
		port:      port,
		creds:     real_creds,
		threadNum: threadsNum,
		timeout:   timeout,
	}
}

func (this *mySQLScan) ScanTarget() []MysqlScanRs {
	ips := common.IPs_SliceOutOfOrder(this.ips)
	p := pool.NewLimited(this.threadNum)
	rs := []MysqlScanRs{}
	defer p.Close()
	batch := p.Batch()

	go func() {
		for i := 0; i < len(ips); i++ {
			ip := ips[i]
			batch.Queue(this.scanHost(ip))

		}
		batch.QueueComplete()
	}()

	for Result := range batch.Results() {

		if err := Result.Error(); err != nil {
			log.Printf("Opps.Something gone wrong!!!%s", err.Error())
			continue
		}
		//fmt.Println("1 port found: ", Result.Value().(int16))
		if len(Result.Value().(MysqlScanRs).Creds) != 0 {
			rs = append(rs, Result.Value().(MysqlScanRs))
		}

	}
	return rs
}

func (this *mySQLScan) scanHost(ip net.IP) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}
		rs := MysqlScanRs{}
		rs.IP = ip
		rs.Creds = nil
		for _, i := range this.creds {

			db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?timeout=%v", i.username, i.password, ip.String(),
				this.port, time.Duration(this.timeout)*time.Second))
			if err != nil {
				log.Printf("[%s] connect error: %e\n", ip.String(), err.Error())
				continue
			}
			cred := credential{}
			cred.username = i.username
			cred.password = i.password
			cred.priv.file = false
			cred.priv.isroot = false
			defer db.Close()
			rootPassRows, err := db.Query("select password from mysql.user where user='root' limit 1")
			if err != nil {
				log.Printf("[%s] query error: %e\n", ip.String(), err.Error())
				continue
			}
			defer rootPassRows.Close()
			var password string
			rootPassRows.Next()
			rootPassRows.Scan(&password)
			if password != "" {
				cred.priv.isroot = true
			}
			filePrivRows, err := db.Query(fmt.Sprintf(
				"select File_priv from mysql.user where user='%s' and File_priv='Y' limit 1", i.username))
			if err != nil {
				log.Printf("[%s] query error: %e\n", ip.String(), err.Error())
				continue
			}
			defer filePrivRows.Close()
			var filePriv string
			filePrivRows.Next()
			filePrivRows.Scan(&filePriv)
			if filePriv != "" {
				cred.priv.file = true
			}

			rs.Creds = append(rs.Creds, cred)
		}

		return rs, nil
	}
}
