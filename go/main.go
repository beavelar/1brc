package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
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

	v1()

	elapsed := time.Since(start)
	fmt.Printf("Took %s to run\n", elapsed)
}

// Super basic tracking and parsing, first go hacking something together
// Average time 1minute 55seconds
// runtime.mapaccess2_faststr 39seconds
// strings.Split 26seconds
// strconv.ParseFloat 14seconds
// runtime.mapassign_faststr 12seconds
// runtime.mapaccess1_faststr 11seconds
// bufio.(*Scanner).Scan 8seconds
// (*Scanner).Text 7.09seconds
func v1() {
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
	for key, _ := range minVals {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)

	output := "{"
	for idx, key := range keys {
		minVal := math.Round(minVals[key] * 10) / 10
		meanVal := math.Round(meanVals[key]/float64(meanCount[key]) * 10) / 10
		maxVal := math.Round(maxVals[key] * 10) / 10
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
