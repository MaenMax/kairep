package uuid

import (
	"fmt"
	"net"
	"os"
	"sync"

	"empowerthings.com/autoendpoint/utils"
)

const (
	FUUID_NB_OF_BYTE = 12
)

var (
	_base        []byte
	_base_length int
	_counter     uint64
	_protect     sync.Mutex
)

type Fuuid [FUUID_NB_OF_BYTE]byte

func NewFastUuid() Fuuid {
	var result Fuuid
	var id uint64
	var l int

	_protect.Lock()
	_counter++
	id = _counter
	_protect.Unlock()

	//	result=make([]byte,FUUID_NB_OF_BYTE,FUUID_NB_OF_BYTE)

	copy(result[l:], _base)
	l += _base_length

	copy(result[l:], utils.UnsafeCastUInt64ToBytes(id))

	return result
}

func init() {
	var ip4 net.IP
	var found_ip4 bool = false

	// Used to keep track of where to copy what.
	var l int

	// Used to contain the current Process ID
	var pid16 uint16

	_base_length = 0

	pid := os.Getpid()

	// The default configuration of Linux comes with
	// /proc/sys/kernel/pid_max = 32768 to be backward compatible
	// with older version of kernel. Only 64 bits can actually go
	// up to 2^22 (4 millions) if properly configured.
	// We suppose we are on a non modified system.
	// So a PID can basically hold on a 16 bits by default.
	pid16 = uint16(pid & 0xFFFF)

	//fmt.Fprintf(os.Stderr,"pid: %v, pid16: %v\n",pid,pid16)

	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip4 = ipnet.IP
					found_ip4 = true
					// fmt.Fprintf(os.Stderr,"%v\n",ip4)
					break
				}

			}
		}
	}

	if found_ip4 {
		_base_length = 4
		_base = make([]byte, _base_length, _base_length)

		//fmt.Fprintf(os.Stderr,"Len(ip4)=%v\n",len(ip4))

		n := len(ip4) - 2

		copy(_base[l:], ip4[n:])
		l += len(ip4) - n

		copy(_base[l:], utils.UnsafeCastUInt16ToBytes(pid16))
		l += 2

	} else {
		fmt.Fprintf(os.Stderr, "Failed to find an IPv4 address. Generated UUID won't be unique among different hosts of a same network!\n", ip4)

		_base_length = 4
		_base = make([]byte, _base_length, _base_length)
		_base[0] = 0
		_base[1] = 0
		l = 2
		copy(_base[l:], utils.UnsafeCastUInt16ToBytes(pid16))
		l += 2
	}
}
