package utils

import (
	"digger/common"
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/go-playground/pool.v3"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type sshScan struct {
	ips       []net.IP
	port      int16
	threadNum uint
	timeout   int
	creds     []credential
}

func InitSSHScan(ips []net.IP, port int16, creds []string, threadNum uint, timeout int) sshScan {
	real_creds := []credential{}

	for _, i := range creds {

		username := strings.Split(i, ":")[0]
		password := strings.Split(i, ":")[1]
		real_creds = append(real_creds, InitCred(username, password))
	}
	return sshScan{
		ips:       ips,
		port:      port,
		creds:     real_creds,
		threadNum: threadNum,
		timeout:   timeout,
	}
}

func (this *sshScan) ScanTarget() []SshScanRs {
	ips := common.IPs_SliceOutOfOrder(this.ips)
	p := pool.NewLimited(this.threadNum)
	rs := []SshScanRs{}
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
		if Result.Value().(SshScanRs).Creds != nil {
			rs = append(rs, Result.Value().(SshScanRs))
		}

	}
	return rs
}

func (this *sshScan) scanHost(ip net.IP) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}
		rs := SshScanRs{}
		rs.IP = ip
		rs.Creds = nil
		for _, i := range this.creds {
			config := &ssh.ClientConfig{
				User: i.username,
				Auth: []ssh.AuthMethod{
					ssh.Password(i.password),
				},
				//Config: ssh.Config{
				//	Ciphers: []string{"aes128-cbc"},
				//},
				HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
					return nil
				},
				Timeout: time.Duration(time.Duration(this.timeout) * time.Second),
			}
			//config.Config.Ciphers = append(config.Config.Ciphers, "aes128-cbc")
			_, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip.String(), this.port), config)
			if err == nil {
				rs.Creds = append(rs.Creds, i)
			}
		}

		return rs, nil
	}
}

func (this *sshScan) checkError(err error, info string) {
	if err != nil {
		fmt.Printf("%s. error: %s\n", info, err)
		os.Exit(1)
	}
}
func (this *sshScan) MuxShell(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 3)
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()

	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			t += n
			result := string(buf[:t])
			if strings.Contains(result, "Username:") ||
				strings.Contains(result, "Password:") ||
				strings.Contains(result, "#") {
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}
