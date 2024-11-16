package algorithms

import (
	"crypto/aes"
	"math"
)

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"math/big"
// )

func Encrypt(key []byte, tweak []byte, X []byte, radix uint64) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{0}, err
	}

	// Step 1
	t := uint64(len(tweak))
	n := uint64(len(X))
	u := n / 2
	v := n - u

	// Step 2
	A, B := X[:u], X[u:]

	// Step 3
	b := CeilingDiv(uint64(math.Ceil(float64(v)*math.Log2(float64(radix)))), 8)

	// Step 4
	d := 4*CeilingDiv(b, 4) + 4

	// Step 5
	P := STRmRadix(1, 2, 8)
	P = append(P, STRmRadix(2, 2, 8)...)
	P = append(P, STRmRadix(1, 2, 8)...)
	P = append(P, STRmRadix(radix, 2, 3*8)...)
	P = append(P, STRmRadix(10, 2, 8)...)
	P = append(P, STRmRadix(Mod(u, 256), 2, 8)...)
	P = append(P, STRmRadix(n, 2, 4*8)...)
	P = append(P, STRmRadix(t, 2, 4*8)...)
	// Step 6
	for i := 0; i < 10; i++ {
		// Step 6.i
		Q := tweak
		Q = append(Q, STRmRadix(0, 2, 8*ModInt(0-int64(t)-int64(b)-1, 16))...)
		Q = append(Q, STRmRadix(uint64(i), 2, 8)...)
		Q = append(Q, STRmRadix(NUMradix(B, radix), 2, int64(b*8))...)

		R := append(P, Q...)

		// Step 6.ii
		R, err = PRF(key, R)
		if err != nil {
			return []byte{0}, err
		}

		// Step 6.iii
		S := R
		for j := uint64(1); j < FloorDiv(d, 16); j++ {
			RxorJ, err := XORBytes(R, STRmRadix(j, 2, 16*8))
			if err != nil {
				return []byte{0}, err
			}
			var encryptedBlock []byte
			block.Encrypt(encryptedBlock, RxorJ)
			S = append(S, encryptedBlock...)
		}

		// Step 6.iv
		y := NUM(S)

		// Step 6.v
		m := int64(v)
		if i%2 == 0 {
			m = int64(u)
		}

		// Step 6.vi
		c := NUMradix(A, radix) + y

		// Step 6.vii
		C := STRmRadix(c, radix, m)

		// Step 6.viii
		A = B

		// Step 6.ix
		B = C
	}
	// Step 7
	A = append(A, B...)
	return A, nil
}
