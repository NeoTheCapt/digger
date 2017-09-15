package common

import (
	"math/rand"
	"net"
	"time"
)

func Int16_SliceOutOfOrder(src []int16) []int16 {
	dest := make([]int16, len(src))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

func IPs_SliceOutOfOrder(src []net.IP) []net.IP {
	dest := make([]net.IP, len(src))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}
