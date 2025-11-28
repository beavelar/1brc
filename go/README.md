# 1 Billion Row Challenge (1BRC) - Go

Being the first language I attempted this challenge with, it was interesing profiling the runs between the Windows machine and the Mac machine. At the beginning there was a large gap between the Windows and Mac machines, but at the end once we concurrently processed the file the Mac machine was faster. It was also interesting to see how different the bottle necks were, on the Windows side there was large bottlelnecks in some of the runtime calls, runtime.mapaccess2_faststr for example, this bottleneck was not in the Mac side however when we started optimizing and only getting a second of improvement. For future runs the Windows timing will be used to view improvements but the Mac timings will still be captured

Unsure if this will translate to the other programming languages but it was interesting to see how "small" adjustments could provide a huge improvement, for example tracking values in 1 map instead of 3 maps, determining the position of the semi-colon and doing string slicing instead of using strings.Split, utiling scanner.Bytes instead of scanner.Text, and utiling multiple goroutine works to process the lines each produced substantial time improvements

## Build and Profile

```bash
go build -o 1brc
./1brc
go tool pprof -http=":8080" 1brc cpu.prof
```

## Versions

### V1
Super basic tracking and parsing, first go hacking something together. Starting with a average time of 1 minute 49 seconds where the major bottlenecks where in accessing multiple maps utilizing strings as the map keys, using strings.Split to split the string by the semi-colon seperator, and using strconv.ParseFloat to convert the string number into a float

#### Timings

```
Average time 1minute 49seconds
runtime.mapaccess2_faststr 37seconds
strings.Split 26seconds
strconv.ParseFloat 16seconds
runtime.mapassign_faststr 11seconds
runtime.mapaccess1_faststr 10seconds
(*Scanner).Text 9seconds
bufio.(*Scanner).Scan 8seconds
```

### V2
Reducing the number of maps used from 3 to 1 but otherwise identical to `V1`, removing the runtime.mapassign_faststr and runtime.mapaccess1_faststr bottlenecks, decreasing the average time by ~30 seconds, now at 1 minute 7 seconds

#### Timings

```
Average time 1minute 7seconds
strings.Split 24seconds
runtime.mapaccess2_faststr 16seconds
strconv.ParseFloat 14seconds
bufio(*Scanner).Scan 9seconds
bufio(*Scanner).Text 7seconds
```

#### V1 Snippet

```go
minVals := make(map[string]float64)
meanVals := make(map[string]float64)
meanCount := make(map[string]int)
maxVals := make(map[string]float64)
```

#### V2 Snippet

```go
type Values struct {
	Max   float64
	Min   float64
	Sum   float64
	Count int
}
values := make(map[string]*Values)
```

### V3
Pretty identical to V2 but uses strings.Index and string slicing instead of strings.Split for seperating the key and the value from the line being read, removing the strings.Split bottleneck, decreasing the average time by ~20 seconds, now at 47 seconds. This does introduce a new bottleneck of strings.Index but this bottleneck is less than strings.Split

#### Timings

```
Average time 47seconds
runtime.mapaccess2_faststr 16seconds
strconv.ParseFloat 13seconds
bufio.(*Scanner).Scan 8seconds
bufio.(*Scanner).Text 8 seconds
strings.Index 5seconds
```

#### V2 Snippet

```go
parts := strings.Split(scanner.Text(), ";")
key := parts[0]
var64, err := strconv.ParseFloat(parts[1], 64)
```

#### V3 Snippet

```go
valStr := scanner.Text()
idx := strings.Index(valStr, ";")
key := valStr[:idx]
var64, err := strconv.ParseFloat(valStr[idx+1:], 64)
```

### V4
Mostly the same as V3, but using scanner.Bytes instead of scanner.Text to read in the next line. This did change how we pulled out the key and value from the line as instead of using strings.Index to locate the semi-colon, we loop through the byte slice. For the value, instead of using strconv.ParseFloat to determine the float value, we are now looping through the value byte slice to determine if the value is positive or negative and building the integer part and the decimal part, then combining the parts to determine the float value. This removes the strconv.ParseFloat and bufio.(*Scanner).Text bottlenecks, decreasing the average time by ~10 seconds, now at 38seconds

#### Timings

```
Average 38seconds
runtime.mapaccess2_faststr 14seconds
runtime.slicebytetostring 9seconds
bufio.(*Scanner).Scan 7seconds
```

#### V3 Snippet

```go
valStr := scanner.Text()
idx := strings.Index(valStr, ";")
key := valStr[:idx]
var64, err := strconv.ParseFloat(valStr[idx+1:], 64)
```

#### V4 Snippet

```go
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
```

### V5
Identical to V4 but sets the size of the values map to a initial size of 1,000. This does not seem to have had much of an impact as the average time only decreased by ~1 second, now at 37 seconds.

#### Timings
```
Average 37seconds
runtime.mapaccess2_faststr 16seconds
runtime.slicebytetostring 7seconds
bufio.(*Scanner).Scan 7seconds
```

#### V4 Snippet

```go
values := make(map[string]*Values)
```

#### V5 Snippet

```go
values := make(map[string]*Values, 1000)
```

### V6
Identical to V5 but updates the Values struct to work with int32/int64 instead of float64. The decimal calculation is then done at the very end when the output string is being generated. This does not seem to have had much of an impact as the average time only decreased by ~1 second, now at 36 seconds

#### Timings

```
Average 36seconds
runtime.mapaccess2_faststr 19seconds
runtime.slicebytetostring 7seconds
bufio.(*Scanner).Scan 7seconds
```

#### V5 Snippet

```go
type Values struct {
	Max   float64
	Min   float64
	Sum   float64
	Count int
}
```

#### V6 Snippet

```go
type ValuesV2 struct {
	Min   int32
	Max   int32
	Sum   int64
	Count int32
}
```

### V7
Identical to V6 but utilizing bytes.IndexByte to located the index of the semi-colon in the line instead of manually looping through the bytes. This does not seem to have had much of an impact as the average time only decreased by ~1 second, now at 35 seconds

#### Timings

```
Average 35seconds
runtime.mapaccess2_faststr 18seconds
bufio.(*Scanner).Scan 7seconds
runtime.slicebytetostring 6seconds
```

#### V6 Snippet

```go
idx := -1
for i, b := range lineBytes {
  if b == ';' {
    idx = i
    break
  }
}
```

#### V7 Snippet

```go
idx := bytes.IndexByte(lineBytes, ';')
```

### V8
Starts with V7 as the base, but introduces goroutines to process the lines concurrently and overriding the scanner Split() method to provide chunks of multiple lines per scan. To avoid any bottleneck with heavy RWMutex usage, the goroutine workers will all work with their own maps, and at the end after all lines have been processed, the output of all goroutines will be combined for the final result. Along with this, instead of processing line by line, the goroutines are provided chunks with multiple lines to work with. This decreased average time by ~20 seconds, now at 14 seconds. Since each goroutine is now processing chunks of lines and not a line at a time, a new bottleneck is introduced because of the usage of strings.SplitSeq

#### Timings

```
Average 14seconds
main.V8.func1.SplitSeq.splitSeq.1 29seconds
bufio.(*Scanner).Scan 18seconds
runtime.mcall 6seconds
```

#### V7 Snippet

```go
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
```

#### V8 Snippet

```go
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

// Create chunks of 1000 lines instead of reading line by line, mess around with line
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
```

### V9
Uses V7 as the base instead of V8 to further optimize single thread performance. Instead of using the string of the city for the map key, use a basic hasher to turn the byte slice from the line to a int64 which will be used as the key for the map. The Values struct is also updated to store the city name as well to be able to easily retrieve the city name for the final output. From V8 to V9 the time does increase since we went back to single threaded processing, from 14 seconds to 33 seconds, but from V7 to V9 the time does decrease a little bit by ~2 seconds, now at 33 seconds

#### Timings

```
Average 33seconds
bufio.(*Scanner).Scan 14seconds
runtime.mapaccess2_fast64 12seconds
```

#### V7 Snippet

```go
type ValuesV2 struct {
	Min   int32
	Max   int32
	Sum   int64
	Count int32
}

// ...

key := string(keyBytes)
```

#### V9 Snippet

```go
type ValuesV3 struct {
	City  string
	Min   int32
	Max   int32
	Sum   int64
	Count int32
}

// ...

hasher.Write(keyBytes)
key := int64(hasher.Sum64())
hasher.Reset()

// ...

values[key] = &ValuesV3{City: string(keyBytes), Min: var32, Sum: var64, Max: var32}

// ...

output += fmt.Sprintf("%s=%.1f/%.1f/%.1f", value.City, minVal, meanVal, maxVal)
```

### V10
Combines the improvements from V8 and V9, running the processing over multiple workers and using a int64 as the map key instead of a string, decreasing the time by ~17 seconds (V9 vs V10), now at 16 seconds

#### Timings

```
Average 16seconds
main.V10.func1.SplitSeq.splitSeq.1 28seconds
bufio.(*Scanner).Scan 16seconds
runtime.mcall 5seconds
runtime.gcBgMarkWorker 5seconds
```

### V11
Identical to V10 but updating the number of workers to spin upu to be 1 less than the number of threads on the CPU and increasing the scanner buffer size, decreasing time by ~6 seconds, now at 10 seconds

#### Timings

```
Average 10seconds
main.V11.func1.SplitSeq.splitSeq.1 34seconds
bufio.(*Scanner).Scan 9seconds
gcBgMarkWorker 6seconds
runtime.mcall 3seconds
```

#### V10 Snippet

```go
workers := 10
```

#### V11 Snippet

```go
workers := runtime.NumCPU() - 1

// ...

buf := make([]byte, 0, 64*1024)
scanner.Buffer(buf, 1024*1024)
```