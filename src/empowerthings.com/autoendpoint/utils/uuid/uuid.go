package uuid

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"net"
	"os"
	"sync"
	"time"
	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/utils"
	"empowerthings.com/autoendpoint/utils/encoding"
)

var (
	_seed         []byte
	_seed_length  int
	_uuid_counter int64
	_uuid_mu      sync.Mutex
)

const (
	BYTES_IN_INT32 = 4
	BYTES_IN_INT64 = 8
	SALT           = "ewklFJ;KVC[\\/JZS;KLFJEW;ARF/OEWPTIWa"
)

func init() {

	_uuid_counter = time.Now().UnixNano()

	hostname, err := os.Hostname()

	if err != nil {
		hostname = err.Error()
	}

	pid := os.Getpid()
	cwd, err := os.Getwd()
	if err != nil {
		// If no current working directory is available, then taking the
		// error message instead. We just need some entropy ... so everything
		// is good to take ...
		cwd = err.Error()
	}

	var l int
	_seed = make([]byte, 2000, 2000)

	random := make([]byte, 128)
	_, err = rand.Read(random)
	if err != nil {
		copy(_seed[l:], random)
		l += 128
	}

	stime := time.Now().UnixNano()
	copy(_seed[l:], utils.UnsafeCastInt64ToBytes(stime))
	l += 8

	copy(_seed[l:], utils.UnsafeCastInt32ToBytes(int32(pid)))
	l += 4

	copy(_seed[l:], hostname)
	l += len(hostname)

	copy(_seed[l:], cwd)
	l += len(cwd)

	copy(_seed[l:], SALT)
	l += len(SALT)

	if ifaces, err := net.Interfaces(); err == nil {
		for _, i := range ifaces {

			if addrs, err := i.Addrs(); err == nil {

				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					copy(_seed[l:], ip.String())
					l += len(ip.String())
				}
			}
		}
	}

	hash := sha1.New()
	io.Copy(hash, bytes.NewReader(_seed))
	_seed = hash.Sum(nil)
	l = hash.Size()

	_seed = _seed[0:l]
	_seed_length = l
}

func FromSeed(seed string) string {
	var result string
	var buffer []byte
	var bresult [20]byte
	var l int

	buffer = make([]byte, 200, 200)

	copy(buffer[l:], seed)
	l += len(seed)
	copy(buffer[l:], SALT)
	l += len(SALT)
	buffer = buffer[0:l]

	bresult = sha1.Sum(buffer)

	// Using our ABase64 implementation
	result = encoding.RawStdEncoding.EncodeToString(bresult[0:20])

	// Returning only 20 characters out of the 27 generated.
	// This means that the total ID space size is:
	// 64^20 = 1.329227996x10^36  (instead of the possible 64^27)
	return result[0:20]
}

func NewUuid() string {
	var result string
	var buffer []byte
	var bresult [20]byte
	var l int

	_uuid_mu.Lock()
	cur_time := time.Now().UnixNano()

	// RS - 5/22/2016 - Added a counter because on Windows, the granularity of
	// the time is milliseconds thus creating several exactly same UUID when
	// called several time within the same millisecond. It is also potentially
	// a problem on Linux in case of requests made in parallel in multiple
	// threads.
	_uuid_counter++

	// 36 = 8 bytes for counter + 8 bytes for UnixNano + 20 bytes from _seed
	// Don't make it calculate dynamically else you will
	// lose 10% of performance!
	buffer = make([]byte, 36, 36)

	copy(buffer[l:], utils.UnsafeCastInt64ToBytes(_uuid_counter))
	l += 8

	copy(buffer[l:], utils.UnsafeCastInt64ToBytes(cur_time))
	l += 8

	copy(buffer[l:], _seed)
	l += _seed_length

	bresult = sha1.Sum(buffer)

	// Using our ABase64 implementation
	result = encoding.RawStdEncoding.EncodeToString(bresult[0:20])
	_uuid_mu.Unlock()

	// Returning only 20 characters out of the 27 generated.
	// This means that the total ID space size is:
	// 64^20 = 1.329227996x10^36  (instead of the possible 64^27)
	return result[0:20]
}
func FormatId(chid string,debug bool) string {

	/*
	        This method will format extracted channelID from CassandraDB to be UUAID standard.

		Example:

		Input: 282c8b9184044db89f942c5972d0e55a

		Output: 282c8b91-8404-4db8-9f94-2c5972d0e55a
	*/
	if debug {

		l4g.Info("Formating channelID :%v", chid)

	}
	var formated_chid string
	var first_part string
	for pos, char := range chid {
		first_part = first_part + string(char)
		if pos == 7 || pos == 11 || pos == 15 || pos == 19 {

			first_part = first_part + "-"
			formated_chid = formated_chid + first_part
			first_part = ""
		}
		if pos == 31 {

			formated_chid = formated_chid + first_part
		}
	}
	return formated_chid
}

