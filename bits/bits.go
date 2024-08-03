package bits

func Count8(x uint8) int {
	return count(x)
}

func Count16(x uint16) int {
	return count(uint8(x&0xFF)) + count(uint8(x>>8))
}

func Count32(x uint32) int {
	return count(uint8(x&0xFF)) + count(uint8(x>>8&0xFF)) +
		count(uint8(x>>16&0xFF)) + count(uint8(x>>24))
}

func Count64(x uint64) int {
	return count(uint8(x&0xFF)) + count(uint8(x>>8&0xFF)) +
		count(uint8(x>>16&0xFF)) + count(uint8(x>>24&0xFF)) +
		count(uint8(x>>32&0xFF)) + count(uint8(x>>40&0xFF)) +
		count(uint8(x>>48&0xFF)) + count(uint8(x>>56))
}

func Normalize8(x uint8) uint8 {
	return normalizerArray8[x]
}

func Normalize4(x uint8) uint8 {
	return normalizerArray4[x]
}
