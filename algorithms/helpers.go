package algorithms

func UintToBytes(num uint) []byte {
	var bytes []byte
	for num > 0 {
		bytes = append([]byte{byte(num % 256)}, bytes...) // prepend bytes
		num /= 256
	}
	return bytes
}
