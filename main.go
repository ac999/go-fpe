package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/ac999/go-fpe/algorithms"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// var alphabets = map[string]string{
// 	"base10": "0123456789",
// 	"base26": "abcdefghijklmnopqrstuvwxyz",
// 	"base36": "0123456789abcdefghijklmnopqrstuvwxyz",
// }

func encryptAES(key, plaintext []byte, iv []byte) ([]byte, error) {
	// AES CBC encryption with padding (assuming the input length is padded)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	// Pad plaintext to blockSize if necessary (PKCS7 padding)
	padding := blockSize - len(plaintext)%blockSize
	padText := make([]byte, padding)
	plaintext = append(plaintext, padText...)
	// Pad IV to blockSize if necessary (PKCS7 padding)
	padding = blockSize - len(iv)%blockSize
	padText = make([]byte, padding)
	iv = append(iv, padText...)
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// AES-CBC Decryption with padding (same as before)
func decryptAES(key, ciphertext []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	// Pad IV to blockSize if necessary (PKCS7 padding)
	padding := blockSize - len(iv)%blockSize
	padText := make([]byte, padding)
	iv = append(iv, padText...)
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove padding
	padding = int(plaintext[len(plaintext)-1])
	plaintext = plaintext[:len(plaintext)-padding]

	return plaintext, nil
}

func FF1encrypt(key, tweak, X []byte, radix uint64) ([]byte, error) {
	// Your custom encryption function (replace with your actual code)
	return algorithms.Encrypt(key, tweak, X, radix)
}

func FF1decrypt(key, tweak, X []byte, radix uint64) ([]byte, error) {
	// Your custom encryption function (replace with your actual code)
	return algorithms.Decrypt(key, tweak, X, radix)
}

// Benchmark Encryption and Decryption for both AES-CBC and Custom
func benchmarkEncryptionsDecryption(numRuns int) ([]float64, []float64, []float64, []float64) {
	// Example input data
	keyHex := "2B7E151628AED2A6ABF7158809CF4F3CEF4359D8D580AA4F"
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		log.Fatal(err)
	}

	tweak := []byte{55, 55, 55, 55, 112, 113, 114, 115, 55, 55, 55}
	X := []byte("0123456789abcdefghi")
	radix := uint64(10)

	// Collect encryption and decryption times for AES and custom encryption
	var aesEncTimes, customEncTimes, aesDecTimes, customDecTimes []float64

	for i := 0; i < numRuns; i++ {
		// Measure AES encryption time
		start := time.Now()
		_, err = encryptAES(key, X, tweak)
		if err != nil {
			log.Fatal(err)
		}
		aesEncDuration := time.Since(start).Seconds() * 1000 // milliseconds
		aesEncTimes = append(aesEncTimes, aesEncDuration)

		// Measure custom encryption time
		start = time.Now()
		_, err = FF1encrypt(key, tweak, X, radix)
		if err != nil {
			log.Fatal(err)
		}
		customEncDuration := time.Since(start).Seconds() * 1000 // milliseconds
		customEncTimes = append(customEncTimes, customEncDuration)

		// Measure AES decryption time
		encData, _ := encryptAES(key, X, tweak)
		start = time.Now()
		_, err = decryptAES(key, encData, tweak)
		if err != nil {
			log.Fatal(err)
		}
		aesDecDuration := time.Since(start).Seconds() * 1000000000 // nanoseconds
		aesDecTimes = append(aesDecTimes, aesDecDuration)

		// Measure custom decryption time
		encData, _ = FF1encrypt(key, tweak, X, radix)
		start = time.Now()
		_, err = FF1decrypt(key, tweak, encData, radix)
		if err != nil {
			log.Fatal(err)
		}
		customDecDuration := time.Since(start).Seconds() * 1000000000 // nanoseconds
		customDecTimes = append(customDecTimes, customDecDuration)
	}

	return aesEncTimes, customEncTimes, aesDecTimes, customDecTimes
}

// Save Times to CSV
func saveTimesToCSV(aesEncTimes, customEncTimes, aesDecTimes, customDecTimes []float64, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Test Name", "FF1 Encryption Time (ns)", "AES-CBC Encryption Time (ns)", "FF1 Decryption Time (ns)", "AES-CBC Decryption Time (ns)"})
	if err != nil {
		log.Fatal(err)
	}

	for i := range aesEncTimes {
		err = writer.Write([]string{
			fmt.Sprintf("Test %d", i+1),
			fmt.Sprintf("%.4f", customEncTimes[i]),
			fmt.Sprintf("%.4f", aesEncTimes[i]),
			fmt.Sprintf("%.4f", customDecTimes[i]),
			fmt.Sprintf("%.4f", aesDecTimes[i]),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Plot the comparison between AES-CBC and FF1 encryption and decryption times
func plotComparison(aesEncTimes, customEncTimes, aesDecTimes, customDecTimes []float64) {
	// Prepare data for the plot
	aesEncData := make(plotter.XYs, len(aesEncTimes))
	customEncData := make(plotter.XYs, len(customEncTimes))
	aesDecData := make(plotter.XYs, len(aesDecTimes))
	customDecData := make(plotter.XYs, len(customDecTimes))

	// Populate data for plotting
	for i := 0; i < len(aesEncTimes); i++ {
		aesEncData[i].X = float64(i)
		aesEncData[i].Y = aesEncTimes[i]
		customEncData[i].X = float64(i)
		customEncData[i].Y = customEncTimes[i]
	}
	for i := 0; i < len(aesDecTimes); i++ {
		aesDecData[i].X = float64(i)
		aesDecData[i].Y = aesDecTimes[i]
		customDecData[i].X = float64(i)
		customDecData[i].Y = customDecTimes[i]
	}

	// Create a new plot
	p := plot.New()

	// Set plot labels and title
	p.Title.Text = "Comparison of AES-CBC and FF1 Encryption & Decryption Times"
	p.X.Label.Text = "Test Number"
	p.Y.Label.Text = "Time (ns)"

	// Create line plots for AES-CBC encryption and decryption
	aesEncLine, err := plotter.NewLine(aesEncData)
	if err != nil {
		log.Fatal(err)
	}
	aesEncLine.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red for AES-CBC Encryption
	aesEncLine.LineStyle.Width = vg.Points(1.5)
	aesEncLine.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	aesDecLine, err := plotter.NewLine(aesDecData)
	if err != nil {
		log.Fatal(err)
	}
	aesDecLine.Color = color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange for AES-CBC Decryption
	aesDecLine.LineStyle.Width = vg.Points(1.5)

	// Create line plots for FF1 encryption and decryption
	customEncLine, err := plotter.NewLine(customEncData)
	if err != nil {
		log.Fatal(err)
	}
	customEncLine.Color = color.RGBA{B: 255, A: 255} // Blue for FF1 Encryption
	customEncLine.LineStyle.Width = vg.Points(1.5)
	customEncLine.LineStyle.Dashes = []vg.Length{vg.Points(3), vg.Points(3)}

	customDecLine, err := plotter.NewLine(customDecData)
	if err != nil {
		log.Fatal(err)
	}
	customDecLine.Color = color.RGBA{G: 255, B: 128, A: 255} // Green for FF1 Decryption
	customDecLine.LineStyle.Width = vg.Points(1.5)

	// Add the lines to the plot
	p.Add(aesEncLine, aesDecLine, customEncLine, customDecLine)

	// Add a legend
	p.Legend.Add("AES-CBC Encryption", aesEncLine)
	p.Legend.Add("AES-CBC Decryption", aesDecLine)
	p.Legend.Add("FF1 Encryption", customEncLine)
	p.Legend.Add("FF1 Decryption", customDecLine)
	p.Legend.Top = true

	// Save the plot to a PNG file
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "encryption_decryption_comparison.png"); err != nil {
		log.Fatal(err)
	}
}

func mean(times []float64) float64 {
	var sum float64
	for _, time := range times {
		sum += time
	}
	return sum / float64(len(times))
}

func median(times []float64) float64 {
	sort.Float64s(times)
	mid := len(times) / 2
	if len(times)%2 == 0 {
		return (times[mid-1] + times[mid]) / 2
	}
	return times[mid]
}

func standardDeviation(times []float64, mean float64) float64 {
	var sunsquares float64
	for _, time := range times {
		diff := time - mean
		sunsquares += diff * diff
	}
	return math.Sqrt(sunsquares / float64(len(times)))
}

// Calculate Statistics and Compare Times for Encryption and Decryption
func compareTimes(aesEncTimes, customEncTimes, aesDecTimes, customDecTimes []float64) {
	// Calculate encryption stats
	aesEncMean := mean(aesEncTimes)
	customEncMean := mean(customEncTimes)
	aesEncMedian := median(aesEncTimes)
	customEncMedian := median(customEncTimes)
	aesEncStdDev := standardDeviation(aesEncTimes, aesEncMean)
	customEncStdDev := standardDeviation(customEncTimes, customEncMean)

	// Calculate decryption stats
	aesDecMean := mean(aesDecTimes)
	customDecMean := mean(customDecTimes)
	aesDecMedian := median(aesDecTimes)
	customDecMedian := median(customDecTimes)
	aesDecStdDev := standardDeviation(aesDecTimes, aesDecMean)
	customDecStdDev := standardDeviation(customDecTimes, customDecMean)

	// Print encryption stats
	fmt.Println("AES-CBC Encryption Stats:")
	fmt.Printf("  Mean: %.4f ns\n", aesEncMean)
	fmt.Printf("  Median: %.4f ns\n", aesEncMedian)
	fmt.Printf("  Standard Deviation: %.4f ns\n", aesEncStdDev)

	fmt.Println("\nFF1 Encryption Stats:")
	fmt.Printf("  Mean: %.4f ns\n", customEncMean)
	fmt.Printf("  Median: %.4f ns\n", customEncMedian)
	fmt.Printf("  Standard Deviation: %.4f ns\n", customEncStdDev)

	// Print decryption stats
	fmt.Println("\nAES-CBC Decryption Stats:")
	fmt.Printf("  Mean: %.4f ns\n", aesDecMean)
	fmt.Printf("  Median: %.4f ns\n", aesDecMedian)
	fmt.Printf("  Standard Deviation: %.4f ns\n", aesDecStdDev)

	fmt.Println("\nFF1 Decryption Stats:")
	fmt.Printf("  Mean: %.4f ns\n", customDecMean)
	fmt.Printf("  Median: %.4f ns\n", customDecMedian)
	fmt.Printf("  Standard Deviation: %.4f ns\n", customDecStdDev)

	// Speed comparisons (encryption and decryption)
	fmt.Printf("\nEncryption Speed Comparison (AES-CBC / FF1): %.4f\n", aesEncMean/customEncMean)
	fmt.Printf("Decryption Speed Comparison (AES-CBC / FF1): %.4f\n", aesDecMean/customDecMean)
}

func main() {
	// Benchmark the encryption and decryption
	numRuns := 100000
	aesEncTimes, customEncTimes, aesDecTimes, customDecTimes := benchmarkEncryptionsDecryption(numRuns)

	// Save times to CSV
	saveTimesToCSV(aesEncTimes, customEncTimes, aesDecTimes, customDecTimes, "encryption_decryption_times.csv")

	// Plot the comparison
	plotComparison(aesEncTimes, customEncTimes, aesDecTimes, customDecTimes)

	// Compare encryption and decryption times
	fmt.Printf("Number of runs: %v", numRuns)
	compareTimes(aesEncTimes, customEncTimes, aesDecTimes, customDecTimes)
}
