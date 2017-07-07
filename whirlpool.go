package whirlpool

const (
	cBlock1     = 8
	cBlock2     = 256
	digestBits  = 512
	digestBytes = digestBits >> 3
	rounds      = 10
)

var (
	sBox = []byte{
		0x18, 0x23, 0xc6, 0xE8, 0x87, 0xB8, 0x01, 0x4F, 0x36, 0xA6, 0xd2, 0xF5, 0x79, 0x6F, 0x91, 0x52,
		0x60, 0xBc, 0x9B, 0x8E, 0xA3, 0x0c, 0x7B, 0x35, 0x1d, 0xE0, 0xd7, 0xc2, 0x2E, 0x4B, 0xFE, 0x57,
		0x15, 0x77, 0x37, 0xE5, 0x9F, 0xF0, 0x4A, 0xdA, 0x58, 0xc9, 0x29, 0x0A, 0xB1, 0xA0, 0x6B, 0x85,
		0xBd, 0x5d, 0x10, 0xF4, 0xcB, 0x3E, 0x05, 0x67, 0xE4, 0x27, 0x41, 0x8B, 0xA7, 0x7d, 0x95, 0xd8,
		0xFB, 0xEE, 0x7c, 0x66, 0xdd, 0x17, 0x47, 0x9E, 0xcA, 0x2d, 0xBF, 0x07, 0xAd, 0x5A, 0x83, 0x33,
		0x63, 0x02, 0xAA, 0x71, 0xc8, 0x19, 0x49, 0xd9, 0xF2, 0xE3, 0x5B, 0x88, 0x9A, 0x26, 0x32, 0xB0,
		0xE9, 0x0F, 0xd5, 0x80, 0xBE, 0xcd, 0x34, 0x48, 0xFF, 0x7A, 0x90, 0x5F, 0x20, 0x68, 0x1A, 0xAE,
		0xB4, 0x54, 0x93, 0x22, 0x64, 0xF1, 0x73, 0x12, 0x40, 0x08, 0xc3, 0xEc, 0xdB, 0xA1, 0x8d, 0x3d,
		0x97, 0x00, 0xcF, 0x2B, 0x76, 0x82, 0xd6, 0x1B, 0xB5, 0xAF, 0x6A, 0x50, 0x45, 0xF3, 0x30, 0xEF,
		0x3F, 0x55, 0xA2, 0xEA, 0x65, 0xBA, 0x2F, 0xc0, 0xdE, 0x1c, 0xFd, 0x4d, 0x92, 0x75, 0x06, 0x8A,
		0xB2, 0xE6, 0x0E, 0x1F, 0x62, 0xd4, 0xA8, 0x96, 0xF9, 0xc5, 0x25, 0x59, 0x84, 0x72, 0x39, 0x4c,
		0x5E, 0x78, 0x38, 0x8c, 0xd1, 0xA5, 0xE2, 0x61, 0xB3, 0x21, 0x9c, 0x1E, 0x43, 0xc7, 0xFc, 0x04,
		0x51, 0x99, 0x6d, 0x0d, 0xFA, 0xdF, 0x7E, 0x24, 0x3B, 0xAB, 0xcE, 0x11, 0x8F, 0x4E, 0xB7, 0xEB,
		0x3c, 0x81, 0x94, 0xF7, 0xB9, 0x13, 0x2c, 0xd3, 0xE7, 0x6E, 0xc4, 0x03, 0x56, 0x44, 0x7F, 0xA9,
		0x2A, 0xBB, 0xc1, 0x53, 0xdc, 0x0B, 0x9d, 0x6c, 0x31, 0x74, 0xF6, 0x46, 0xAc, 0x89, 0x14, 0xE1,
		0x16, 0x3A, 0x69, 0x09, 0x70, 0xB6, 0xd0, 0xEd, 0xcc, 0x42, 0x98, 0xA4, 0x28, 0x5c, 0xF8, 0x86}

	c  = make([][]uint64, cBlock1)
	rc = make([]uint64, rounds)
)

func init() {

	for i := range c {
		c[i] = make([]uint64, cBlock2)
	}

	for x := range c[0] {
		v1 := uint64(sBox[x])

		v2 := v1 << 1

		if v2 >= 0x100 {
			v2 ^= 0x11d
		}

		v4 := v2 << 1

		if v4 >= 0x100 {
			v4 ^= 0x11d
		}

		v5 := v4 ^ v1
		v8 := v4 << 1

		if v8 >= 0x100 {
			v8 ^= 0x11d
		}

		v9 := v8 ^ v1

		c[0][x] = (v1 << 56) | (v1 << 48) | (v4 << 40) | (v1 << 32) | (v8 << 24) | (v5 << 16) | (v2 << 8) | (v9)

		for t := 1; t < cBlock1; t++ {
			c[t][x] = (c[t-1][x] >> 8) | (c[t-1][x] << 56)
		}

	}

	for k := range rc {
		i := 8 * (k)
		rc[k] =
			(c[0][i] & 0xff00000000000000) ^
				(c[1][i+1] & 0x00ff000000000000) ^
				(c[2][i+2] & 0x0000ff0000000000) ^
				(c[3][i+3] & 0x000000ff00000000) ^
				(c[4][i+4] & 0x00000000ff000000) ^
				(c[5][i+5] & 0x0000000000ff0000) ^
				(c[6][i+6] & 0x000000000000ff00) ^
				(c[7][i+7] & 0x00000000000000ff)
	}

}

type nessieStruct struct {
	bitLength  []byte
	buffer     []byte
	bufferBits uint32
	bufferPos  uint32
	hash       []uint64
	k          []uint64
	l          []uint64
	block      []uint64
	state      []uint64
}

func makeNessie() *nessieStruct {
	return &nessieStruct{
		make([]byte, 32),
		make([]byte, 64),
		0, 0,
		make([]uint64, 8),
		make([]uint64, 8),
		make([]uint64, 8),
		make([]uint64, 8),
		make([]uint64, 8)}
}

func (nessie *nessieStruct) processBuffer() {
	j := 0
	for i := range nessie.block {
		v := uint64(nessie.buffer[j])
		j++
		v = (v << 8) ^ uint64(nessie.buffer[j])
		j++
		v = (v << 8) ^ uint64(nessie.buffer[j])
		j++
		v = (v << 8) ^ uint64(nessie.buffer[j])
		j++
		v = (v << 8) ^ uint64(nessie.buffer[j])
		j++
		v = (v << 8) ^ uint64(nessie.buffer[j])
		j++
		v = (v << 8) ^ uint64(nessie.buffer[j])
		j++
		v = (v << 8) ^ uint64(nessie.buffer[j])
		j++

		nessie.block[i] = v
	}

	for i := range nessie.block {
		nessie.k[i] = nessie.hash[i]
		nessie.state[i] = nessie.block[i] ^ nessie.k[i]
	}

	for _, v := range rc {
		for i := range nessie.block {
			nessie.l[i] = 0

			s := uint(56)
			for t := 0; t < 8; t++ {
				nessie.l[i] ^= c[t][(nessie.k[(i-t)&7]>>s)&0xff]
				s -= 8
			}
		}

		copy(nessie.k[:8], nessie.l[:8])

		nessie.k[0] ^= v
		for i := range nessie.block {
			nessie.l[i] = nessie.k[i]

			s := uint(56)
			for t := 0; t < 8; t++ {
				nessie.l[i] ^= c[t][(nessie.state[(i-t)&7]>>s)&0xff]
				s -= 8
			}
		}

		copy(nessie.state[:8], nessie.l[:8])
	}

	for i := range nessie.block {
		nessie.hash[i] ^= nessie.state[i] ^ nessie.block[i]
	}
}

func (nessie *nessieStruct) add(source []byte) {

	sourceBits := uint32(len(source) * 8)
	sourcePos := uint32(0)
	sourceGap := uint32((8 - (sourceBits & 7)) & 7)
	bufferRem := uint32(nessie.bufferBits & 7)

	var b uint8

	value := sourceBits
	carry := uint32(0)

	for i := 31; i >= 0; i-- {
		carry += uint32((nessie.bitLength[i] & 0xff)) + (value & 0xff)
		nessie.bitLength[i] = byte(carry)
		carry >>= 8
		value >>= 8
	}

	for sourceBits > 8 {
		b = ((source[sourcePos] << sourceGap) & 0xff) | ((source[sourcePos+1] & 0xff) >> (8 - sourceGap))

		nessie.buffer[nessie.bufferPos] |= b >> bufferRem
		nessie.bufferPos++
		nessie.bufferBits += 8 - bufferRem
		if nessie.bufferBits == 512 {
			nessie.processBuffer()
			nessie.bufferBits = 0
			nessie.bufferPos = 0
		}
		nessie.buffer[nessie.bufferPos] = ((b << (8 - bufferRem)) & 0xff)
		nessie.bufferBits += bufferRem
		sourceBits -= 8
		sourcePos++
	}

	if sourceBits > 0 {
		b = (source[sourcePos] << sourceGap) & 0xff
		nessie.buffer[nessie.bufferPos] |= b >> bufferRem
	} else {
		b = 0
	}

	if bufferRem+sourceBits < 8 {
		nessie.bufferBits += sourceBits
	} else {
		nessie.bufferPos++
		nessie.bufferBits += 8 - bufferRem
		sourceBits -= 8 - bufferRem
		if nessie.bufferBits == 512 {
			nessie.processBuffer()
			nessie.bufferBits = 0
			nessie.bufferPos = 0
		}
		nessie.buffer[nessie.bufferPos] = ((b << (8 - bufferRem)) & 0xff)
		nessie.bufferBits += sourceBits
	}

}

func (nessie *nessieStruct) finalize() []byte {

	nessie.buffer[nessie.bufferPos] |= 0x80 >> (nessie.bufferBits & 7)
	nessie.bufferPos++
	if nessie.bufferPos > 32 {
		for nessie.bufferPos < 64 {
			nessie.buffer[nessie.bufferPos] = 0
			nessie.bufferPos++
		}
		nessie.processBuffer()
		nessie.bufferPos = 0
	}
	for nessie.bufferPos < 32 {
		nessie.buffer[nessie.bufferPos] = 0
		nessie.bufferPos++
	}

	copy(nessie.buffer[32:64], nessie.bitLength[:32])

	nessie.processBuffer()

	digest := make([]byte, digestBytes)
	i := 0
	j := 0
	for i < digestBytes/8 {
		h := nessie.hash[i]
		digest[j] = byte(h >> 56)
		j++
		digest[j] = byte(h >> 48)
		j++
		digest[j] = byte(h >> 40)
		j++
		digest[j] = byte(h >> 32)
		j++
		digest[j] = byte(h >> 24)
		j++
		digest[j] = byte(h >> 16)
		j++
		digest[j] = byte(h >> 8)
		j++
		digest[j] = byte(h)
		j++
		i++
	}

	return digest

}

// Whirlpool struct.
// Do not instantiate this yourself. Use InitWhirlpool() instead.
type Whirlpool struct {
	nessie *nessieStruct
}

// InitWhirlpool is called to start using the Whirlpool hash algorithm.
func InitWhirlpool() *Whirlpool {
	w := &Whirlpool{makeNessie()}
	return w
}

// Write data to the Whirlpool buffer
func (wp *Whirlpool) Write(p []byte) (n int, err error) {
	wp.nessie.add(p)
	return len(p), nil
}

// Digest returns a finished hash in the form of a byte array
func (wp *Whirlpool) Digest() []byte {
	return wp.nessie.finalize()
}
