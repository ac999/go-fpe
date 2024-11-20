package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/ac999/go-fpe/algorithms"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
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

func FF1encrypt(key, tweak, X []byte, radix uint64) ([]byte, error) {
	// Your custom encryption function (replace with your actual code)
	return algorithms.Encrypt(key, tweak, X, radix)
}

func benchmarkEncryptions(numRuns int) ([]float64, []float64) {
	// Example input data
	keyHex := "2B7E151628AED2A6ABF7158809CF4F3CEF4359D8D580AA4F"
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		log.Fatal(err)
	}

	tweak := []byte{55, 55, 55, 55, 112, 113, 114, 115, 55, 55, 55}
	X := []byte("0123456789abcdefghi")
	radix := uint64(10)

	// Collect encryption times for AES and custom encryption
	var aesTimes []float64
	var customTimes []float64

	for i := 0; i < numRuns; i++ {
		// Measure AES encryption time
		start := time.Now()
		_, err = encryptAES(key, X, tweak)
		if err != nil {
			log.Fatal(err)
		}
		aesDuration := time.Since(start).Seconds() * 1000000000 // nanoseconds
		aesTimes = append(aesTimes, aesDuration)

		// Measure custom encryption time
		start = time.Now()
		_, err = FF1encrypt(key, tweak, X, radix)
		if err != nil {
			log.Fatal(err)
		}
		customDuration := time.Since(start).Seconds() * 1000000000 // nanoeconds
		customTimes = append(customTimes, customDuration)
	}

	return aesTimes, customTimes
}

func saveTimesToCSV(aesTimes, customTimes []float64, filename string) {
	// Create or open the CSV file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	err = writer.Write([]string{"Test Name", "FF1 Encryption Time (ns)", "AES-CBC Encryption Time (ns)"})
	if err != nil {
		log.Fatal(err)
	}

	// Write the times to CSV
	for i := range aesTimes {
		err = writer.Write([]string{
			fmt.Sprintf("Test %d", i+1),
			fmt.Sprintf("%.4f", customTimes[i]), // FF1 encryption time
			fmt.Sprintf("%.4f", aesTimes[i]),    // AES encryption time
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func plotComparison(aesTimes, customTimes []float64) {
	p := plot.New()

	// Create the plot data
	aesPoints := make(plotter.XYs, len(aesTimes))
	customPoints := make(plotter.XYs, len(customTimes))

	for i := range aesTimes {
		aesPoints[i].X = float64(i)
		aesPoints[i].Y = aesTimes[i]
		customPoints[i].X = float64(i)
		customPoints[i].Y = customTimes[i]
	}

	// Plot AES data
	aesLine, err := plotter.NewLine(aesPoints)
	if err != nil {
		log.Fatal(err)
	}
	aesLine.Color = plotutil.Color(0)

	// Plot custom encryption data
	customLine, err := plotter.NewLine(customPoints)
	if err != nil {
		log.Fatal(err)
	}
	customLine.Color = plotutil.Color(1)

	// Add lines to the plot
	p.Add(aesLine, customLine)

	// Set labels and title
	p.Title.Text = "Encryption Speed Comparison"
	p.X.Label.Text = "Run Number"
	p.Y.Label.Text = "Time (ns)"

	// Save plot as a PNG file
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "comparison_plot.png"); err != nil {
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
	var sumSquares float64
	for _, time := range times {
		diff := time - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(times)))
}

// func loadCSV(filename string) ([]float64, []float64, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer file.Close()

// 	reader := csv.NewReader(file)
// 	// Skip header
// 	_, err = reader.Read()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	var aesTimes, customTimes []float64
// 	for {
// 		record, err := reader.Read()
// 		if err != nil {
// 			break
// 		}

// 		// Parse the times from CSV (assuming they are in the second and third columns)
// 		customTime, err := strconv.ParseFloat(record[1], 64)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		aesTime, err := strconv.ParseFloat(record[2], 64)
// 		if err != nil {
// 			return nil, nil, err
// 		}

// 		customTimes = append(customTimes, customTime)
// 		aesTimes = append(aesTimes, aesTime)
// 	}

// 	return aesTimes, customTimes, nil
// }

func compareTimes(aesTimes, customTimes []float64) {
	// Calculate stats for AES-CBC
	aesMean := mean(aesTimes)
	aesMedian := median(aesTimes)
	aesStdDev := standardDeviation(aesTimes, aesMean)

	// Calculate stats for FF1 (custom encryption)
	customMean := mean(customTimes)
	customMedian := median(customTimes)
	customStdDev := standardDeviation(customTimes, customMean)

	// Print statistics
	fmt.Println("AES-CBC Encryption Stats:")
	fmt.Printf("  Mean: %.4f ns\n", aesMean)
	fmt.Printf("  Median: %.4f ns\n", aesMedian)
	fmt.Printf("  Standard Deviation: %.4f ns\n", aesStdDev)

	fmt.Println("\nFF1 Encryption Stats:")
	fmt.Printf("  Mean: %.4f ns\n", customMean)
	fmt.Printf("  Median: %.4f ns\n", customMedian)
	fmt.Printf("  Standard Deviation: %.4f ns\n", customStdDev)

	// Calculate speed comparison (speed factor)
	speedFactor := aesMean / customMean
	fmt.Printf("\nSpeed Comparison (AES-CBC / FF1): %.4f\n", speedFactor)
}

func main() {
	// Define the number of runs for each encryption method
	numRuns := 100000

	// Benchmark both encryption methods
	aesTimes, customTimes := benchmarkEncryptions(numRuns)

	// Save the encryption times to CSV
	saveTimesToCSV(aesTimes, customTimes, "encryption_times.csv")

	// Plot the results
	plotComparison(aesTimes, customTimes)

	// Compare the times and print the statistics
	compareTimes(aesTimes, customTimes)
}
