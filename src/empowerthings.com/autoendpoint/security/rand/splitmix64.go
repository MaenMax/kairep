package rand


type SplitMix64 struct {
	_x int64
}

func NewSplitMix64(x int64) *SplitMix64 {

	var tmp SplitMix64=SplitMix64{}

	tmp._x=x
	
	return &tmp
}

func (x *SplitMix64) Next() int64 {
	//x._x+=0x9E3779B97F4A7C15
	x._x+=-0x61C8864680B583EB
	z:=(x._x)
	//z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 30)) * -0x40A7B892E31B1A47
	//z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	z = (z ^ (z >> 27)) * -0x6B2FB644ECCEEE15
	return z ^ (z >> 31);
}

