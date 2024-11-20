package algorithms

import (
	"crypto/aes"
	"math"
	"math/big"
)

func Encrypt(key []byte, tweak []byte, X []byte, radix uint64) ([]byte, error) {
	BigRadix := big.NewInt(int64(radix))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Step 1
	t := uint64(len(tweak))
	n := uint64(len(X))
	u := n / 2
	v := n - u

	// Step 2
	A, B := X[:u], X[u:]

	// Step 3
	// Calculate log2(radix):
	log2Radix := math.Log2(float64(radix))

	b := CeilingDiv(uint64(math.Ceil(float64(v)*log2Radix)), 8)

	// Step 4
	d := 4*CeilingDiv(b, 4) + 4

	// Step 5
	P := []byte{1, 2, 1}
	P = append(P, BigSTRmRadix(BigRadix, 256, 3)...)
	P = append(P, BigSTRmRadix(big.NewInt(10), 256, 1)...)
	P = append(P, BigSTRmRadix(new(big.Int).SetUint64(Mod(u, 256)), 256, 1)...)
	P = append(P, BigSTRmRadix(new(big.Int).SetUint64(n), 256, 4)...)
	P = append(P, BigSTRmRadix(new(big.Int).SetUint64(t), 256, 4)...)

	// Step 6
	for i := int64(0); i < 10; i++ {

		// Step 6.i
		Q := tweak
		Q = append(Q, BigSTRmRadix(big.NewInt(0), 256, ModInt(16-int64(t)-int64(b)-1, 16))...)
		Q = append(Q, BigSTRmRadix(big.NewInt(i), 256, 1)...)
		Q = append(Q, BigSTRmRadix(BigNUMradix(B, radix), 256, int64(b))...)

		R := P
		R = append(R, Q...)

		// Step 6.ii
		R, err := PRF(key, R)
		if err != nil {
			return []byte{0}, err
		}

		// Step 6.iii
		S := R

		for j := uint64(1); j < CeilingDiv(d, 16); j++ {
			RxorJ, err := XORBytes(R, BigSTRmRadix(new(big.Int).SetUint64(j), 256, 16))

			if err != nil {
				return []byte{0}, err
			}
			encryptedBlock := make([]byte, aes.BlockSize)
			block.Encrypt(encryptedBlock, RxorJ)

			S = append(S, encryptedBlock...)

		}
		S = S[:d]

		// Step 6.iv
		y := BigNUM(S)

		// Step 6.v
		mBig := new(big.Int).SetUint64(v)
		m := int64(v)
		if i%2 == 0 {
			mBig = new(big.Int).SetUint64(u)
			m = int64(u)
		}

		// Step 6.vi
		BigAplusY := BigNUMradix(A, radix)
		BigAplusY = BigAplusY.Add(BigAplusY, y)
		radixAtM := BigPower(BigRadix, mBig)
		c := BigMod(BigAplusY, radixAtM)

		// Step 6.vii
		C := BigSTRmRadix(c, radix, m)

		// Step 6.viii
		A = B

		// Step 6.ix
		B = C

	}
	// Step 7
	A = append(A, B...)
	return A, nil
}

func Decrypt(key []byte, tweak []byte, X []byte, radix uint64) ([]byte, error) {
	BigRadix := big.NewInt(int64(radix))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Step 1
	t := uint64(len(tweak))
	n := uint64(len(X))
	u := n / 2
	v := n - u

	// Step 2
	A, B := X[:u], X[u:]

	// Step 3
	// Calculate log2(radix):
	log2Radix := math.Log2(float64(radix))

	b := CeilingDiv(uint64(math.Ceil(float64(v)*log2Radix)), 8)

	// Step 4
	d := 4*CeilingDiv(b, 4) + 4

	// Step 5
	P := []byte{1, 2, 1}
	P = append(P, BigSTRmRadix(BigRadix, 256, 3)...)
	P = append(P, BigSTRmRadix(big.NewInt(10), 256, 1)...)
	P = append(P, BigSTRmRadix(new(big.Int).SetUint64(Mod(u, 256)), 256, 1)...)
	P = append(P, BigSTRmRadix(new(big.Int).SetUint64(n), 256, 4)...)
	P = append(P, BigSTRmRadix(new(big.Int).SetUint64(t), 256, 4)...)

	// Step 6
	for i := int64(9); i >= 0; i-- {

		// Step 6.i
		Q := tweak
		Q = append(Q, BigSTRmRadix(big.NewInt(0), 256, ModInt(0-int64(t)-int64(b)-1, 16))...)
		Q = append(Q, BigSTRmRadix(big.NewInt(i), 256, 1)...)
		Q = append(Q, BigSTRmRadix(BigNUMradix(A, radix), 256, int64(b))...)

		R := P
		R = append(R, Q...)

		// Step 6.ii
		R, err := PRF(key, R)
		if err != nil {
			return []byte{0}, err
		}

		// Step 6.iii
		S := R

		for j := uint64(1); j < CeilingDiv(d, 16); j++ {
			RxorJ, err := XORBytes(R, BigSTRmRadix(new(big.Int).SetUint64(j), 256, 16))

			if err != nil {
				return []byte{0}, err
			}
			encryptedBlock := make([]byte, aes.BlockSize)
			block.Encrypt(encryptedBlock, RxorJ)
			S = append(S, encryptedBlock...)
		}
		S = S[:d]

		// Step 6.iv
		y := BigNUM(S)

		// Step 6.v
		mBig := new(big.Int).SetUint64(v)
		m := int64(v)
		if i%2 == 0 {
			mBig = new(big.Int).SetUint64(u)
			m = int64(u)
		}

		// Step 6.vi

		BigBminusY := BigNUMradix(B, radix)
		BigBminusY = BigBminusY.Sub(BigBminusY, y)
		c := BigMod(BigBminusY, BigPower(BigRadix, mBig))

		// Step 6.vii
		C := BigSTRmRadix(c, radix, m)

		// Step 6.viii
		B = A

		// Step 6.ix
		A = C

	}
	// Step 7
	A = append(A, B...)
	return A, nil
}
