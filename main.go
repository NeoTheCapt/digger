package main

import (
	"digger/common"
	"digger/utils"
	"fmt"
	"github.com/alexflint/go-arg"
	"log"
	"os"
)

type args struct {
	MODE      string   `arg:"required,help:scan or mysql or ssh."`
	IP        string   `arg:"required,help:The IPs which you want to scan"`
	PORT      []int16  `arg:"help:The ports which you want to scan"`
	Creds     []string `arg:"help:MySql credentials."` //123:321 234:123
	THREADNUM uint     `arg:"help:number of threads"`
	WORDERNUM uint     `arg:"help:number of workers"`
	TIMEOUT   int      `arg:"help:seconds of timeout"`
}

func (args) Description() string {
	return "e.g: \n" +
		"digger --ip 10.0.0.1-10.0.0.254" +
		""
}
func main() {
	var args args
	//set default value
	args.THREADNUM = 2
	args.WORDERNUM = 2
	args.TIMEOUT = 2
	args.Creds = []string{"root:root", "root:123456"}
	arg.MustParse(&args)
	IPs, err := common.ParseIPStr(args.IP)
	if err != nil {
		log.Printf("error found in ParseIPStr: %e\n", err.Error())
		os.Exit(1)
	}
	if len(args.PORT) == 0 {
		log.Println("no specific ports")
		os.Exit(1)
	}
	switch args.MODE {
	case "scan":
		var rs []utils.PortScanRs
		ps := utils.PortScan{IPs: IPs, Ports: args.PORT,
			ThreadNum: args.THREADNUM, WorkersNum: args.WORDERNUM,
			Timeout: args.TIMEOUT,
		}
		rs = ps.ScanTargets()
		println("Done:")
		for i := 0; i < len(rs); i++ {
			fmt.Printf("IP :%s , Open ports: %v\n", rs[i].IP.String(), rs[i].Ports)
		}
		break
	case "mysql":
		log.Println("Start mysql scan")
		var rs []utils.MysqlScanRs
		mysqlScan := utils.InitMySQLScan(IPs, args.PORT[0], args.Creds, args.THREADNUM, args.TIMEOUT)
		rs = mysqlScan.ScanTarget()
		for _, i := range rs {
			println(fmt.Sprintf("[%s] login with [%v]", i.IP, i.Creds))
		}
		break
	case "ssh":
		log.Println("Start ssh scan")
		var rs []utils.SshScanRs
		sshScan := utils.InitSSHScan(IPs, args.PORT[0], args.Creds, args.THREADNUM, args.TIMEOUT)
		rs = sshScan.ScanTarget()
		for _, i := range rs {
			println(fmt.Sprintf("[%s] login with [%v]", i.IP, i.Creds))
		}
		break
	}
}
