package algorithms

import (
	"crypto/aes"
	"fmt"
	"math"
)

func Encrypt(key []byte, tweak []byte, X []byte, radix uint64) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	fmt.Printf("##########\n")
	fmt.Printf("Key is %v\nRadix = %v\nPT is <%v>\nTweak is <%v>\n", key, radix, X, tweak)
	// Step 1
	t := uint64(len(tweak))
	n := uint64(len(X))
	u := n / 2
	v := n - u
	fmt.Printf("\nStep 1: u is %v, v is %v\n", u, v)

	// Step 2
	A, B := X[:u], X[u:]
	fmt.Printf("Step 2: A is %v\nB is %v\n", A, B)

	// Step 3
	b := CeilingDiv(uint64(math.Ceil(float64(v)*math.Log2(float64(radix)))), 8)
	fmt.Printf("Step 3: b is %v\n", b)

	// Step 4
	d := 4*CeilingDiv(b, 4) + 4
	fmt.Printf("Step 4: d is %v\n", d)

	// Step 5
	P := []byte{1, 2, 1}
	P = append(P, STRmRadix(radix, 256, 3)...)
	P = append(P, STRmRadix(10, 256, 1)...)
	P = append(P, STRmRadix(Mod(u, 256), 256, 1)...)
	P = append(P, STRmRadix(n, 256, 4)...)
	P = append(P, STRmRadix(t, 256, 4)...)
	fmt.Printf("Step 5: P is %v\n", P)
	// Step 6
	for i := 0; i < 10; i++ {
		fmt.Printf("\nRound #%v\n", i)
		// Step 6.i
		Q := tweak
		Q = append(Q, STRmRadix(0, 256, ModInt(0-int64(t)-int64(b)-1, 16))...)
		Q = append(Q, STRmRadix(uint64(i), 256, 1)...)
		Q = append(Q, STRmRadix(NUMradix(B, radix), 256, int64(b))...)

		fmt.Printf("Step 6.i Q is %v\n", Q)

		R := P
		R = append(R, Q...)

		// Step 6.ii
		R, err := PRF(key, R)
		if err != nil {
			return []byte{0}, err
		}
		fmt.Printf("Step 6.ii R is %v\n", R)

		// Step 6.iii
		var S []byte
		for j := uint64(1); j < CeilingDiv(d, 16); j++ {
			RxorJ, err := XORBytes(R, STRmRadix(j, 256, 16))
			fmt.Printf("Step 6.iii iteration %v Step 1: RxorJ is %v\n", j, RxorJ)
			if err != nil {
				return []byte{0}, err
			}
			encryptedBlock := make([]byte, aes.BlockSize)
			block.Encrypt(encryptedBlock, RxorJ)
			fmt.Printf("Step 6.iii iteration %v Step 2: encryptedBlock is %s\n", j, encryptedBlock)
			fmt.Printf("Step 6.iii iteration %v Step 3: S before append is %s\n", j, S)
			S = append(S, encryptedBlock...)
			fmt.Printf("Step 6.iii iteration %v Step 4: S after append is %s\n", j, S)
		}
		fmt.Printf("Step 6.iii S is %v\n", S)
		// Step 6.iv
		y := NUM(S)
		fmt.Printf("Step 6.iv y is %v\n", y)

		// Step 6.v
		m := int64(v)
		if i%2 == 0 {
			m = int64(u)
		}

		// Step 6.vi
		c := Mod(NUMradix(A, radix)+y, Power(radix, uint64(m)))

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

func Decrypt(key []byte, tweak []byte, X []byte, radix uint64) ([]byte, error) {
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
	b := CeilingDiv(uint64(math.Ceil(float64(v)*math.Log2(float64(radix)))), 8)

	// Step 4
	d := 4*CeilingDiv(b, 4) + 4

	// Step 5
	P := STRmRadix(1, 2, 8)
	P = append(P, STRmRadix(2, 256, 1)...)
	P = append(P, STRmRadix(1, 256, 1)...)
	P = append(P, STRmRadix(radix, 256, 3)...)
	P = append(P, STRmRadix(10, 256, 1)...)
	P = append(P, STRmRadix(Mod(u, 256), 256, 1)...)
	P = append(P, STRmRadix(n, 256, 4)...)
	P = append(P, STRmRadix(t, 256, 4)...)
	// Step 6
	for i := 9; i >= 0; i-- {
		// Step 6.i
		Q := tweak
		Q = append(Q, STRmRadix(0, 256, ModInt(0-int64(t)-int64(b)-1, 16))...)
		Q = append(Q, STRmRadix(uint64(i), 256, 1)...)
		Q = append(Q, STRmRadix(NUMradix(A, radix), 256, int64(b))...)

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
			RxorJ, err := XORBytes(R, STRmRadix(j, 256, 16))
			if err != nil {
				return []byte{0}, err
			}
			encryptedBlock := make([]byte, aes.BlockSize)
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
		c := Mod(NUMradix(B, radix)-y, Power(radix, uint64(m)))

		// Step 6.vii
		C := STRmRadix(c, radix, m)

		// Step 6.viii
		B = A

		// Step 6.ix
		A = C
	}
	// Step 7
	A = append(A, B...)
	return A, nil
}
