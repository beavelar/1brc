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

### V3

Updates the file reading portion to use FileReader and provide the contents to a char array, iterate over that char array instead of reading through the buffered reader

#### Timings

```
Average time 1minute 32seconds
java/io/Reader.read 1,847 samples
java/lang/Double.parseDouble 1,556 samples
java/lang/StringBuilder.append 1,354 samples
java/lang/StringBuilder.toString 1,649 samples
java/util/HashMap.get 1,648 samples
```

#### V2 Snippet

```java
try (var reader = new BufferedReader(new FileReader("../1brc/measurements.txt"), 12 * 1024 * 1024)) {
  String line;
  while ((line = reader.readLine()) != null) {
    var parts = line.split(";");
    var key = parts[0];
    var value = Double.parseDouble(parts[1]);

    //
  }
}
```

#### V2 Snippet

```java
try (var reader = new FileReader("../1brc/measurements.txt")) {
  var citySb = new StringBuilder();
  var tempSb = new StringBuilder();
  var parsingCity = true;
  var charsRead = 0;
  var city = "";
  var temp = 0.0;

  while ((charsRead = reader.read(cbuf)) > 0) {
    for (var idx = 0; idx < charsRead; idx++) {
      if (cbuf[idx] == ';') {
        parsingCity = false;
        city = citySb.toString();
        citySb = new StringBuilder();
        continue;
      }

      if (cbuf[idx] == '\r' || cbuf[idx] == '\n') {
        parsingCity = true;
        temp = Double.parseDouble(tempSb.toString());
        tempSb = new StringBuilder();

        //
        continue;
      }

      if (parsingCity) {
        citySb.append(cbuf[idx]);
      } else {
        tempSb.append(cbuf[idx]);
      }
    }
  }
}
```
