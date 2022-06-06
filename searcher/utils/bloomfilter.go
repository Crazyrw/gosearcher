package utils

import "math"

// Filter is an encoded set of []byte keys.
type Filter []byte

// MayContainKey _
func (f Filter) MayContainKey(k []byte) bool {
	return f.MayContain(Hash(k))
}

// MayContain returns whether the filter may contain given key. False positives
// are possible, where it returns true for keys not in the original set.
func (f Filter) MayContain(h uint32) bool {
	//Implement me here!!!
	if len(f) < 2 {
		return false
	}
	//看 h 是否 包含在 过滤器里面
	// h 是已经经过hash计算的 32位 哈希值

	//获取hash次数
	k := f[len(f)-1]
	if k > 30 {
		return true
	}
	//获取数组长度 多少位
	//因为最后一个 字节是 用来存放哈希次数的
	nBits := (len(f) - 1) * 8
	delta := h>>17 | h<<15

	for i := uint8(0); i < k; i++ {
		temp := h % uint32(nBits)
		if f[temp/8]&(1<<(temp%8)) == 0 {
			return false
		}
		h += delta
	}
	return true
}

// NewFilter returns a new Bloom filter that encodes a set of []byte keys with
// the given number of bits per key, approximately.
//
// A good bitsPerKey value is 10, which yields a filter with ~ 1% false
// positive rate.
//返回的是，插入了一组 key值的布隆过滤器，key的长度为 bitsPerKey
func NewFilter(keys []uint32, bitsPerKey int) Filter {
	return Filter(appendFilter(keys, bitsPerKey))
}

// BloomBitsPerKey returns the bits per key required by bloomfilter based on
// the false positive rate.
func BloomBitsPerKey(numEntries int, fp float64) int {
	//Implement me here!!!
	//阅读bloom论文实现，并在这里编写公式
	//传入参数numEntries是bloom中存储的数据个数，fp是false positive假阳性率
	//即  函数作用是 返回 m/n的值，即每一个key占据数组的大小
	size := -1 * float64(numEntries) * math.Log(fp) / math.Pow(math.Ln2, 2)
	//向上取整
	temp := math.Ceil(size / float64(numEntries))
	return int(temp)
}

func appendFilter(keys []uint32, bitsPerKey int) []byte {
	//Implement me here!!!
	//在这里实现将多个Key值放入到bloom过滤器中
	if bitsPerKey < 0 {
		bitsPerKey = 0
	}
	//计算 哈希次数
	k := uint32(math.Ln2 * float64(bitsPerKey))
	if k < 1 {
		k = 1
	}
	if k > 32 {
		k = 32
	}
	//计算 布隆过滤器大小
	//使用 字节数组来存
	nBits := len(keys) * bitsPerKey
	if nBits < 64 {
		nBits = 64
	}
	//这里 + 7的意思是，比如 7 / 8 = 0, 那需要7位，但是申请不到一个8位的字节
	nBytes := (nBits + 7) / 8
	//重新更新一下 布隆过滤器的大小
	nBits = nBytes * 8
	//最后一个字节用来存放 hash次数 K
	Filter := make([]byte, nBytes+1)
	//插入 keys
	for _, h := range keys {
		delta := h>>17 | h<<15
		for j := uint32(0); j < k; j++ {
			temp := h % uint32(nBits)
			Filter[temp/8] |= 1 << (temp % 8)
			h += delta
		}
	}
	Filter[nBytes] = byte(k)
	return Filter
}

// Hash implements a hashing algorithm similar to the Murmur hash.
func Hash(b []byte) uint32 {
	//Implement me here!!!
	//在这里实现高效的HashFunction
	const (
		seed = 0xbc9f1d34
		m    = 0xc6a4a793
	)
	h := uint32(seed) ^ uint32(len(b))*m
	for ; len(b) >= 4; b = b[4:] {
		h += uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
		h *= m
		h ^= h >> 16
	}
	switch len(b) {
	case 3:
		h += uint32(b[2]) << 16
		fallthrough
	case 2:
		h += uint32(b[1]) << 8
		fallthrough
	case 1:
		h += uint32(b[0])
		h *= m
		h ^= h >> 24
	}
	return h
}
