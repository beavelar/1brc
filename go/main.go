package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("Running calculations")
	start := time.Now()
	prof, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	defer prof.Close()
	pprof.StartCPUProfile(prof)
	defer pprof.StopCPUProfile()

	// V1()
	// V2()
	// V3()
	// V4()
	// V5()
	// V6()
	// V7()
	V8()

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run\n", elapsed)
}

// Super basic tracking and parsing, first go hacking something together
//
// Average time 1minute 49seconds
// runtime.mapaccess2_faststr 37seconds
// strings.Split 26seconds
// strconv.ParseFloat 16seconds
// runtime.mapassign_faststr 11seconds
// runtime.mapaccess1_faststr 10seconds
// (*Scanner).Text 9seconds
// bufio.(*Scanner).Scan 8seconds
//
// Mac Average time 2minute 25seconds
func V1() {
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	minVals := make(map[string]float64)
	meanVals := make(map[string]float64)
	meanCount := make(map[string]int)
	maxVals := make(map[string]float64)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ";")
		key := parts[0]
		var64, err := strconv.ParseFloat(parts[1], 64)

		if err != nil {
			log.Fatal(err)
		}

		// Min eval
		val, found := minVals[key]
		if !found {
			minVals[key] = var64
		} else {
			if val > var64 {
				minVals[key] = var64
			}
		}

		// Mean eval
		val, found = meanVals[key]
		if !found {
			meanVals[key] = var64
		} else {
			meanVals[key] = val + var64
			meanCount[key] = meanCount[key] + 1
		}

		// Max eval
		val, found = maxVals[key]
		if !found {
			maxVals[key] = var64
		} else {
			if val < var64 {
				maxVals[key] = var64
			}
		}
	}

	keys := make([]string, len(minVals))
	idx := 0
	for key := range minVals {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := math.Round(minVals[key]*10) / 10
		meanVal := math.Round(meanVals[key]/float64(meanCount[key])*10) / 10
		maxVal := math.Round(maxVals[key]*10) / 10
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, minVal, meanVal, maxVal)
		if idx < len(keys)-1 {
			output += ", "
		}
	}
	output += "}"
	fmt.Println(output)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type Values struct {
	Max   float64
	Min   float64
	Sum   float64
	Count int
}

// Reducing the number of maps used but mostly the same as v1
//
// Average time 1minute 7seconds
// strings.Split 24seconds
// runtime.mapaccess2_faststr 16seconds
// strconv.ParseFloat 14seconds
// bufio(*Scanner).Scan 9seconds
// bufio(*Scanner).Text 7seconds
//
// Mac Average time 1minute 37seconds
func V2() {
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	values := make(map[string]*Values)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ";")
		key := parts[0]
		var64, err := strconv.ParseFloat(parts[1], 64)

		if err != nil {
			log.Fatal(err)
		}

		val, found := values[key]
		if !found {
			values[key] = &Values{Min: var64, Sum: var64, Max: var64}
		} else {
			// Min eval
			if val.Min > var64 {
				val.Min = var64
			}

			// Mean eval
			val.Sum = val.Sum + var64
			val.Count = val.Count + 1

			// Max eval
			if val.Max < var64 {
				val.Max = var64
			}
		}
	}

	keys := make([]string, len(values))
	idx := 0
	for key := range values {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := math.Round(values[key].Min*10) / 10
		meanVal := math.Round(values[key].Sum/float64(values[key].Count)*10) / 10
		maxVal := math.Round(values[key].Max*10) / 10
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, minVal, meanVal, maxVal)
		if idx < len(keys)-1 {
			output += ", "
		}
	}
	output += "}"
	fmt.Println(output)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Identical to V2 but opts for string slicing instead of using strings.Split
//
// Average time 47seconds
// runtime.mapaccess2_faststr 16seconds
// strconv.ParseFloat 13seconds
// bufio.(*Scanner).Scan 8seconds
// bufio.(*Scanner).Text 8 seconds
// strings.Index 5seconds
//
// Mac Average time 1minute 8seconds
func V3() {
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	values := make(map[string]*Values)
	for scanner.Scan() {
		valStr := scanner.Text()
		idx := strings.Index(valStr, ";")
		key := valStr[:idx]
		var64, err := strconv.ParseFloat(valStr[idx+1:], 64)

		if err != nil {
			log.Fatal(err)
		}

		val, found := values[key]
		if !found {
			values[key] = &Values{Min: var64, Sum: var64, Max: var64}
		} else {
			// Min eval
			if val.Min > var64 {
				val.Min = var64
			}

			// Mean eval
			val.Sum = val.Sum + var64
			val.Count = val.Count + 1

			// Max eval
			if val.Max < var64 {
				val.Max = var64
			}
		}
	}

	keys := make([]string, len(values))
	idx := 0
	for key := range values {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := math.Round(values[key].Min*10) / 10
		meanVal := math.Round(values[key].Sum/float64(values[key].Count)*10) / 10
		maxVal := math.Round(values[key].Max*10) / 10
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, minVal, meanVal, maxVal)
		if idx < len(keys)-1 {
			output += ", "
		}
	}
	output += "}"
	fmt.Println(output)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Mostly identical to V3 but using scanner.Bytes() instead of scanner.Text()
//
// Average 38seconds
// runtime.mapaccess2_faststr 14seconds
// runtime.slicebytetostring 9seconds
// bufio.(*Scanner).Scan 7seconds
//
// Mac Average 57seconds
func V4() {
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	values := make(map[string]*Values)
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		idx := -1
		for i, b := range lineBytes {
			if b == ';' {
				idx = i
				break
			}
		}

		keyBytes := lineBytes[:idx]
		valBytes := lineBytes[idx+1:]
		key := string(keyBytes)

		var sign float64 = 1.0
		var intPart, fracPart int64
		var decimalSeen bool
		var numStart int

		if valBytes[0] == '-' {
			sign = -1.0
			numStart = 1
		} else {
			numStart = 0
		}

		for i := numStart; i < len(valBytes); i++ {
			if valBytes[i] == '.' {
				decimalSeen = true
				continue
			}
			digit := int64(valBytes[i] - '0')
			if !decimalSeen {
				intPart = intPart*10 + digit
			} else {
				fracPart = digit
			}
		}
		var64 := sign * (float64(intPart) + float64(fracPart)/10.0)

		if err != nil {
			log.Fatal(err)
		}

		val, found := values[key]
		if !found {
			values[key] = &Values{Min: var64, Sum: var64, Max: var64}
		} else {
			// Min eval
			if val.Min > var64 {
				val.Min = var64
			}

			// Mean eval
			val.Sum = val.Sum + var64
			val.Count = val.Count + 1

			// Max eval
			if val.Max < var64 {
				val.Max = var64
			}
		}
	}

	keys := make([]string, len(values))
	idx := 0
	for key := range values {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := math.Round(values[key].Min*10) / 10
		meanVal := math.Round(values[key].Sum/float64(values[key].Count)*10) / 10
		maxVal := math.Round(values[key].Max*10) / 10
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, minVal, meanVal, maxVal)
		if idx < len(keys)-1 {
			output += ", "
		}
	}
	output += "}"
	fmt.Println(output)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Pretty much the save as V4 but sets the size of the values slice to 1,000
//
// Average 37seconds
// runtime.mapaccess2_faststr 16seconds
// runtime.slicebytetostring 7seconds
// bufio.(*Scanner).Scan 7seconds
//
// Mac Average 55seconds
func V5() {
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	values := make(map[string]*Values, 1000)
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		idx := -1
		for i, b := range lineBytes {
			if b == ';' {
				idx = i
				break
			}
		}

		keyBytes := lineBytes[:idx]
		valBytes := lineBytes[idx+1:]
		key := string(keyBytes)

		var sign float64 = 1.0
		var intPart, fracPart int64
		var decimalSeen bool
		var numStart int

		if valBytes[0] == '-' {
			sign = -1.0
			numStart = 1
		} else {
			numStart = 0
		}

		for i := numStart; i < len(valBytes); i++ {
			if valBytes[i] == '.' {
				decimalSeen = true
				continue
			}
			digit := int64(valBytes[i] - '0')
			if !decimalSeen {
				intPart = intPart*10 + digit
			} else {
				fracPart = digit
			}
		}
		var64 := sign * (float64(intPart) + float64(fracPart)/10.0)

		if err != nil {
			log.Fatal(err)
		}

		val, found := values[key]
		if !found {
			values[key] = &Values{Min: var64, Sum: var64, Max: var64}
		} else {
			// Min eval
			if val.Min > var64 {
				val.Min = var64
			}

			// Mean eval
			val.Sum = val.Sum + var64
			val.Count = val.Count + 1

			// Max eval
			if val.Max < var64 {
				val.Max = var64
			}
		}
	}

	keys := make([]string, len(values))
	idx := 0
	for key := range values {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := math.Round(values[key].Min*10) / 10
		meanVal := math.Round(values[key].Sum/float64(values[key].Count)*10) / 10
		maxVal := math.Round(values[key].Max*10) / 10
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, minVal, meanVal, maxVal)
		if idx < len(keys)-1 {
			output += ", "
		}
	}
	output += "}"
	fmt.Println(output)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Store the values as int and do the final float calculation at the very end
type ValuesV2 struct {
	Min   int32
	Max   int32
	Sum   int64
	Count int32
}

// Mostly the save as V5 but opts for working with int32 and int64 for the values
// to trace and do the float64 work only at the end
//
// Average 36seconds
// runtime.mapaccess2_faststr 19seconds
// runtime.slicebytetostring 7seconds
// bufio.(*Scanner).Scan 7seconds
//
// Mac Average 54seconds
func V6() {
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	values := make(map[string]*ValuesV2, 1000)
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		idx := -1
		for i, b := range lineBytes {
			if b == ';' {
				idx = i
				break
			}
		}

		keyBytes := lineBytes[:idx]
		valBytes := lineBytes[idx+1:]
		key := string(keyBytes)

		var sign int32 = 1
		var intPart, fracPart int32
		var decimalSeen bool
		var numStart int

		if valBytes[0] == '-' {
			sign = -1
			numStart = 1
		} else {
			numStart = 0
		}

		for i := numStart; i < len(valBytes); i++ {
			if valBytes[i] == '.' {
				decimalSeen = true
				continue
			}
			digit := int32(valBytes[i] - '0')
			if !decimalSeen {
				intPart = intPart*10 + digit
			} else {
				fracPart = digit
			}
		}
		var32 := sign * (intPart*10 + fracPart)
		var64 := int64(var32)

		if err != nil {
			log.Fatal(err)
		}

		if val, found := values[key]; !found {
			values[key] = &ValuesV2{Min: var32, Sum: var64, Max: var32}
		} else {
			// Min eval
			if val.Min > var32 {
				val.Min = var32
			}

			// Mean eval
			val.Sum += var64
			val.Count++

			// Max eval
			if val.Max < var32 {
				val.Max = var32
			}
		}
	}

	keys := make([]string, len(values))
	idx := 0
	for key := range values {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := float64(values[key].Min) / 10
		meanVal := math.Round(float64(values[key].Sum)/float64(values[key].Count)*10) / 100
		maxVal := float64(values[key].Max) / 10
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, minVal, meanVal, maxVal)
		if idx < len(keys)-1 {
			output += ", "
		}
	}
	output += "}"
	fmt.Println(output)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Identical to V6 but utilizing bytes.IndexByte to locate the semicolon instead
// of manually searching using the for loop
//
// Average 35seconds
// runtime.mapaccess2_faststr 18seconds
// bufio.(*Scanner).Scan 7seconds
// runtime.slicebytetostring 6seconds
//
// Mac Average 54seconds
func V7() {
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	values := make(map[string]*ValuesV2, 1000)
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		idx := bytes.IndexByte(lineBytes, ';')

		keyBytes := lineBytes[:idx]
		valBytes := lineBytes[idx+1:]
		key := string(keyBytes)

		var sign int32 = 1
		var intPart, fracPart int32
		var decimalSeen bool
		var numStart int

		if valBytes[0] == '-' {
			sign = -1
			numStart = 1
		} else {
			numStart = 0
		}

		for i := numStart; i < len(valBytes); i++ {
			if valBytes[i] == '.' {
				decimalSeen = true
				continue
			}
			digit := int32(valBytes[i] - '0')
			if !decimalSeen {
				intPart = intPart*10 + digit
			} else {
				fracPart = digit
			}
		}
		var32 := sign * (intPart*10 + fracPart)
		var64 := int64(var32)

		if err != nil {
			log.Fatal(err)
		}

		if val, found := values[key]; !found {
			values[key] = &ValuesV2{Min: var32, Sum: var64, Max: var32}
		} else {
			// Min eval
			if val.Min > var32 {
				val.Min = var32
			}

			// Mean eval
			val.Sum += var64
			val.Count++

			// Max eval
			if val.Max < var32 {
				val.Max = var32
			}
		}
	}

	keys := make([]string, len(values))
	idx := 0
	for key := range values {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := float64(values[key].Min) / 10
		meanVal := math.Round(float64(values[key].Sum)/float64(values[key].Count)*10) / 100
		maxVal := float64(values[key].Max) / 10
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, minVal, meanVal, maxVal)
		if idx < len(keys)-1 {
			output += ", "
		}
	}
	output += "}"
	fmt.Println(output)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Starts with the base of V7 but overrides the scanner Split() method to return a string
// chunk of multiple lines instead of line by line, and spins up multiple workers to process
// the chunk by splitting it by newlines and then running through the same calculations.
// Combines the results of all workers at the end to produce the end result
//
// Mac Average 13seconds
func V8() {
	file, err := os.Open("../1brc/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// The number of workers to spin up to handle line chunk processing/calculations, mess
	// around with the number of workers to view the impact
	workers := 10

	var wg sync.WaitGroup
	linesChan := make(chan string, 10000)
	resultMaps := make([]map[string]*ValuesV2, workers)

	for idx := range workers {
		wg.Add(1)
		resultMap := make(map[string]*ValuesV2)
		resultMaps[idx] = resultMap
		go func(wg *sync.WaitGroup, input chan string, output map[string]*ValuesV2) {
			for chunkStr := range input {
				for lineStr := range strings.SplitSeq(chunkStr, "\n") {
					idx := strings.Index(lineStr, ";")

					keyBytes := lineStr[:idx]
					valBytes := lineStr[idx+1:]
					key := string(keyBytes)

					var sign int32 = 1
					var intPart, fracPart int32
					var decimalSeen bool
					var numStart int

					if valBytes[0] == '-' {
						sign = -1
						numStart = 1
					} else {
						numStart = 0
					}

					for i := numStart; i < len(valBytes); i++ {
						if valBytes[i] == '.' {
							decimalSeen = true
							continue
						}
						digit := int32(valBytes[i] - '0')
						if !decimalSeen {
							intPart = intPart*10 + digit
						} else {
							fracPart = digit
						}
					}
					var32 := sign * (intPart*10 + fracPart)
					var64 := int64(var32)

					if err != nil {
						log.Fatal(err)
					}

					if val, found := output[key]; !found {
						output[key] = &ValuesV2{Min: var32, Sum: var64, Max: var32}
					} else {
						// Min eval
						if val.Min > var32 {
							val.Min = var32
						}

						// Mean eval
						val.Sum += var64
						val.Count++

						// Max eval
						if val.Max < var32 {
							val.Max = var32
						}
					}
				}
			}
			wg.Done()
		}(&wg, linesChan, resultMap)
	}

	scanner := bufio.NewScanner(file)

	// Create chunks of 100 lines instead of reading line by line, mess around with line
	// chunks to view the impact
	linesPerChunk := 1000
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		newlineCount := 0
		lastNewlineIndex := -1
		for i, b := range data {
			if b == '\n' {
				newlineCount++
				lastNewlineIndex = i
			}
			if newlineCount >= linesPerChunk {
				return lastNewlineIndex + 1, data[:lastNewlineIndex], nil
			}
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	values := make(map[string]*ValuesV2, 1000)
	for scanner.Scan() {
		linesChan <- scanner.Text()
	}

	close(linesChan)
	wg.Wait()

	for _, resultMap := range resultMaps {
		for key, val := range resultMap {
			if finalVal, found := values[key]; !found {
				values[key] = val
			} else {
				if finalVal.Min > val.Min {
					finalVal.Min = val.Min
				}
				finalVal.Sum += val.Sum
				finalVal.Count += val.Count
				if finalVal.Max < val.Max {
					finalVal.Max = val.Max
				}
			}
		}
	}

	keys := make([]string, len(values))
	idx := 0
	for key := range values {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := float64(values[key].Min) / 10
		meanVal := math.Round(float64(values[key].Sum)/float64(values[key].Count)*10) / 100
		maxVal := float64(values[key].Max) / 10
		output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", key, minVal, meanVal, maxVal)
		if idx < len(keys)-1 {
			output += ", "
		}
	}
	output += "}"
	fmt.Println(output)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
