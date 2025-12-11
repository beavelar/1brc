# 1 Billion Row Challenge (1BRC) - Java

## Run and Profile

In order to profile, need to have [async-profiler](https://github.com/async-profiler/async-profiler) downloaded locally

**Run on Mac**

```bash
java -agentpath:/Users/brian/files/async-profiler/async-profiler-4.2.1-macos/lib/libasyncProfiler.dylib=start,event=cpu,file=profile.html _1brc.java
```

** Run on Windows**

```bash
java -agentpath:/home/brian/installs/async-profiler/async-profiler-4.2.1-linux-x64/lib/libasyncProfiler.so=start,event=cpu,file=profile.html _1brc.java
```

## Versions

### V1

Just a simple first cut, reading the file by lines and using a HashMap to track the unique citys, sorting the keys after the file has been parsed and printing the outcome

#### Timings

```
Average time 1minute 34seconds
java/io/BufferedReader.readLine 3,545 samples
java/lang/Double.parseDouble 1,402 samples
java/lang/String.split 3,329 samples
java/util/HashMap.get 1,325 samples
```

### V2

Identical to V2 but increases the buffer size of BufferedReader to 12MB instead of the default 8KB. Timing wise also identical to V1

#### Timings

```
Average time 1minute 34seconds
java/io/BufferedReader.readLine 3,740 samples
java/lang/Double.parseDouble 1,226 samples
java/lang/String.split
java/util/HashMap.get 1,314 samples
```

#### V1 Snippet

```java
try (var reader = new BufferedReader(new FileReader("../1brc/measurements.txt"))) {
  //
}
```

#### V2 Snippet

```java
try (var reader = new BufferedReader(new FileReader("../1brc/measurements.txt"), 12 * 1024 * 1024)) {
  //
}
```
