package rand

import (
	"fmt"
)

type Xoroshiro128Plus struct {
	state [2]uint64
}

func NewXoroshiro128Plus(s1 uint64, s2 uint64) *Xoroshiro128Plus {

	var tmp Xoroshiro128Plus=Xoroshiro128Plus{}

	tmp.state[0]=s1
	tmp.state[1]=s2
	
	return &tmp
}


func (x *Xoroshiro128Plus) Next() uint64 {
	var s0, s1 uint64
	s0 = x.state[0]
	s1 = x.state[1]
	result:=s0+s1;

	s1^=s0
	x.state[0] = lro(s0,55) ^ s1 ^ (s1<<14)
	x.state[1] = lro(s1,36)

	return result
}

func lro(v uint64, n uint8) uint64 {

	if n>=64 {
		panic(fmt.Sprint("Invalid parameter n=%v",n));
	}

	return (v<<n) | (v>>(64-n))
}
