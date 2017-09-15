package utils

import (
	"digger/common"
	"fmt"
	"gopkg.in/go-playground/pool.v3"
	"log"
	"net"
	"time"
)

type PortScan struct {
	IPs        []net.IP
	Ports      []int16
	ThreadNum  uint
	WorkersNum uint
	Timeout    int
	ScanType   string //T:Connect scan S:SYN scan

}

func (this *PortScan) ScanTargets() []PortScanRs {
	ips := common.IPs_SliceOutOfOrder(this.IPs)
	p := pool.NewLimited(this.ThreadNum)
	rs := []PortScanRs{}
	defer p.Close()
	batch := p.Batch()

	go func() {
		for i := 0; i < len(ips); i++ {
			ip := ips[i]
			batch.Queue(this.ScanHost(ip))

		}
		batch.QueueComplete()
	}()

	for Result := range batch.Results() {

		if err := Result.Error(); err != nil {
			log.Printf("Opps.Something gone wrong!!!%s", err.Error())
			continue
		}
		//fmt.Println("1 port found: ", Result.Value().(int16))
		if len(Result.Value().(PortScanRs).Ports) != 0 {
			rs = append(rs, Result.Value().(PortScanRs))
		}

	}
	return rs
}
func (this *PortScan) ScanHost(ip net.IP) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}
		rs := PortScanRs{}
		rs = this.scanHost(ip)
		return rs, nil
	}
}

func (this *PortScan) scanHost(ip net.IP) PortScanRs {
	portlist := common.Int16_SliceOutOfOrder(this.Ports)
	p := pool.NewLimited(this.WorkersNum)
	rs := PortScanRs{}
	rs.IP = ip
	defer p.Close()
	batch := p.Batch()

	go func() {
		for i := 0; i < len(portlist); i++ {
			//color.Green(IP)
			batch.Queue(this.scanPort(ip.String(), portlist[i]))
		}

		batch.QueueComplete()
	}()

	for Result := range batch.Results() {

		if err := Result.Error(); err != nil {
			log.Printf("Opps.Something gone wrong!!!%s", err.Error())
			continue
		}
		//fmt.Println("1 port found: ", Result.Value().(int16))
		if Result.Value().(int16) != 0 {
			rs.Ports = append(rs.Ports, Result.Value().(int16))
		}

	}
	return rs
}

func (this *PortScan) scanPort(ip string, port int16) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}
		rs := false
		switch this.ScanType {
		case "T":
			rs = this.connectScan(ip, port)
			break
		case "S":
			//@todo add SYN scan
			break
		default:
			rs = this.connectScan(ip, port)
		}
		if rs {
			return port, nil
		} else {
			return int16(0), nil
		}
	}
}

func (this *PortScan) connectScan(ipStr string, port int16) bool {
	log.Printf("work start on: %s:%d", ipStr, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", ipStr, port))
	if err != nil {
		return false
	}
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Duration(this.Timeout)*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
