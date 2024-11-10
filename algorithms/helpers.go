package algorithms

import "encoding/binary"

func UintToBytes(num uint64) []byte {
	var bytes = make([]byte, 8)
	binary.PutUvarint(bytes, num)
	return bytes
}

func power(x uint64, y uint64) uint64 {
	result := uint64(1)
	for i := uint64(0); i < y; i++ {
		result *= x
	}
	return result
}

func intToNBytes(x int, n int) []byte {

	result := make([]byte, n)
	for i := n - 1; i >= 0; i-- {
		result[i] = byte(x % 10)
		x /= 10
	}
	return []byte(result)
}

func uint64ToNBytes(x uint64, n int) []byte {
	result := make([]byte, n)
	for i := n - 1; i >= 0; i-- {
		result[i] = byte(x % 10)
		x /= 10
	}
	return []byte(result)
}
