package common

import (
	"net"
	"strings"
)

func CIDR2IPs(s string) ([]net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}

	// Get network bits.
	nb, _ := ipnet.Mask.Size()

	// Short-circuit a /32.
	if nb == 32 {
		return []net.IP{ip}, nil
	}

	// Get usable hosts count. Init slice with the size;
	// /16 and larger will start to cause runtime slow downs
	// if we used append().
	nHosts := 2<<uint(31-nb) - 2
	ips := make([]net.IP, nHosts)

	// net.IP slice start position.
	p := 15

	// Increment start IP by # hosts and populate slice.
	for n := 0; n < nHosts; n++ {

		ip[15]++
		// Increment the next class if the current position is > 254.
		if ip[15] > 254 {
			ip[15] = 1
		decPos:
			p--
			for n := 15 - p; n > 0; n-- {
				// Break at Class A limit.
				if n > 3 {
					break
				}
				ip[15-n]++
				if ip[15-n] > 254 {
					ip[15-n] = 0
					goto decPos
				}
			}
			ips[n] = net.ParseIP(ip.String())
		} else {
			ips[n] = net.ParseIP(ip.String())
		}
	}

	return ips, nil
}

func CIDR2IPs2(s string) ([]net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, net.ParseIP(ip.String()))
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Range2IPs(s string) ([]net.IP, error) {
	ips := []net.IP{}
	ipStart := net.ParseIP(strings.Split(s, "-")[0])
	ipEnd := net.ParseIP(strings.Split(s, "-")[1])
	for ip := ipStart; !ipStart.Equal(ipEnd); inc(ip) {
		ips = append(ips, net.ParseIP(ip.String()))
	}
	ips = append(ips, ipEnd)
	return ips, nil
}

func ParseIPStr(IPStr string) ([]net.IP, error) {
	if strings.Contains(IPStr, "/") {
		//log.Println("IP with mask")
		IPs, err := CIDR2IPs(IPStr)
		return IPs, err
	}
	if strings.Contains(IPStr, "-") {
		//log.Println("IP region")
		IPs, err := Range2IPs(IPStr)
		return IPs, err
	}
	if strings.Contains(IPStr, ",") {
		//log.Println("Sigle IP")
		IPs := []net.IP{}
		ipStrArr := strings.Split(IPStr, ",")
		for _, i := range ipStrArr {
			IPs = append(IPs, net.ParseIP(i))
		}
		return IPs, nil
	}
	IPs := []net.IP{net.ParseIP(IPStr)}
	return IPs, nil
}
