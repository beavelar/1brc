# 1 Billion Row Challenge (1BRC)

Attempting the 1 billion row challenge in various languages and profiling to learn more about language specific optimizations

## The Challenge

The 1BRC involves reading a text file with 1 billion rows (each `city;temperature`) and calculating the min, max, and mean temperature for each city, then printing the results sorted by city alphabetically.

[`Link to official challenge details`](https://github.com/gunnarmorling/1brc)

## Setup

1. Clone this repo, navigate to the directory, and clone the [`1 Billion Row Challenge Repo`](https://github.com/gunnarmorling/1brc)

```bash
git clone git@github.com:beavelar/1brc.git
cd 1brc
git clone git@github.com:gunnarmorling/1brc.git
```

2. Navigate to the inner `1brc` directory and build (Requires JDK 21 and Maven)

```bash
cd 1brc
mvn clean verify
```

3. Generate the measurements.txt file

```bash
bash create_measurements.sh 1000000000
```

*All attempts will re-use this measurements.txt file within the inner 1brc directory to not generate unecessary copies of the file as the generated file is around 13.5GB*

## Go

**Build and Profile**

```bash
go build -o 1brc
./1brc
go tool pprof -http=":8080" 1brc cpu.pprof
```

### Versions

#### V1
Super basic tracking and parsing, first go hacking something together. Starting with a average time of 1 minute 49 seconds where the major bottlenecks where in accessing multiple maps utilizing strings as the map keys, using strings.Split to split the string by the semi-colon seperator, and using strconv.ParseFloat to convert the string number into a float

**Timings**

```
runtime.mapaccess2_faststr 37seconds
strings.Split 26seconds
strconv.ParseFloat 16seconds
runtime.mapassign_faststr 11seconds
runtime.mapaccess1_faststr 10seconds
(*Scanner).Text 9seconds
bufio.(*Scanner).Scan 8seconds
```

#### V2
Reducing the number of maps used from 3 to 1 but otherwise identical to `V1`, removing the runtime.mapassign_faststr and runtime.mapaccess1_faststr bottlenecks, decreasing the average time by ~30 seconds, now at 1 minute 7 seconds

**Timings**

```
Average time 1minute 7seconds
strings.Split 24seconds
runtime.mapaccess2_faststr 16seconds
strconv.ParseFloat 14seconds
bufio(*Scanner).Scan 9seconds
bufio(*Scanner).Text 7seconds
```

**V1 Snippet**

```go
minVals := make(map[string]float64)
meanVals := make(map[string]float64)
meanCount := make(map[string]int)
maxVals := make(map[string]float64)
```

**V2 Snippet**

```go
type Values struct {
	Max   float64
	Min   float64
	Sum   float64
	Count int
}
values := make(map[string]*Values)
```

#### V3
Pretty identical to V2 but uses strings.Index and string slicing instead of strings.Split for seperating the key and the value from the line being read, removing the strings.Split bottleneck, decreasing the average time by ~20 seconds, now at 47 seconds. This does introduce a new bottleneck of strings.Index but this bottleneck is less than strings.Split

**Timings**

```
Average time 47seconds
runtime.mapaccess2_faststr 16seconds
strconv.ParseFloat 13seconds
bufio.(*Scanner).Scan 8seconds
bufio.(*Scanner).Text 8 seconds
strings.Index 5seconds
```

**V2 Snippet**

```go
parts := strings.Split(scanner.Text(), ";")
key := parts[0]
var64, err := strconv.ParseFloat(parts[1], 64)
```

**V3 Snippet**

```go
valStr := scanner.Text()
idx := strings.Index(valStr, ";")
key := valStr[:idx]
var64, err := strconv.ParseFloat(valStr[idx+1:], 64)
```

#### V4
Mostly the same as V3, but using scanner.Bytes instead of scanner.Text to read in the next line. This did change how we pulled out the key and value from the line as instead of using strings.Index to locate the semi-colon, we loop through the byte slice. For the value, instead of using strconv.ParseFloat to determine the float value, we are now looping through the value byte slice to determine if the value is positive or negative and building the integer part and the decimal part, then combining the parts to determine the float value. This removes the strconv.ParseFloat and bufio.(*Scanner).Text bottlenecks, decreasing the average time by ~10 seconds, now at 38seconds

**Timings**

```
Average 38seconds
runtime.mapaccess2_faststr 14seconds
runtime.slicebytetostring 9seconds
bufio.(*Scanner).Scan 7seconds
```

**V3 Snippet**

```go
valStr := scanner.Text()
idx := strings.Index(valStr, ";")
key := valStr[:idx]
var64, err := strconv.ParseFloat(valStr[idx+1:], 64)
```

**V4 Snippet**

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

#### V5
Identical to V4 but sets the size of the values map to a initial size of 1,000. This does not seem to have had much of an impact as the average time only decreased by ~1 second, now at 37 seconds.

**Timings**
```
Average 37seconds
runtime.mapaccess2_faststr 16seconds
runtime.slicebytetostring 7seconds
bufio.(*Scanner).Scan 7seconds
```

**V4 Snippet**

```go
values := make(map[string]*Values)
```

**V5 Snippet**

```go
values := make(map[string]*Values, 1000)
```

#### V6
Identical to V5 but updates the Values struct to work with int32/int64 instead of float64. The decimal calculation is then done at the very end when the output string is being generated. This does not seem to have had much of an impact as the average time only decreased by ~1 second, now at 36 seconds

**Timings**

```
Average 36seconds
runtime.mapaccess2_faststr 19seconds
runtime.slicebytetostring 7seconds
bufio.(*Scanner).Scan 7seconds
```

**V5 Snippet**

```go
type Values struct {
	Max   float64
	Min   float64
	Sum   float64
	Count int
}
```

**V6 Snippet**

```go
type ValuesV2 struct {
	Min   int32
	Max   int32
	Sum   int64
	Count int32
}
```

#### V7
Identical to V6 but utilizing bytes.IndexByte to located the index of the semi-colon in the line instead of manually looping through the bytes. This does not seem to have had much of an impact as the average time only decreased by ~1 second, now at 35 seconds

**Timings**

```
Average 35seconds
runtime.mapaccess2_faststr 18seconds
bufio.(*Scanner).Scan 7seconds
runtime.slicebytetostring 6seconds
```

**V6 Snippet**

```go
idx := -1
for i, b := range lineBytes {
  if b == ';' {
    idx = i
    break
  }
}
```

**V7 Snippet**

```go
idx := bytes.IndexByte(lineBytes, ';')
```