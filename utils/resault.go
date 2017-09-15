package utils

import "net"

type MysqlScanRs struct {
	IP    net.IP
	Creds []credential //验证过的帐号
}
type PortScanRs struct {
	IP    net.IP
	Ports []int16 //打开的端口
}
type SshScanRs struct {
	IP    net.IP
	Creds []credential
}
