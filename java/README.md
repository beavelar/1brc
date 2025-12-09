# 1 Billion Row Challenge (1BRC) - Java

## Run and Profile

In order to profile, need to have [async-profiler](https://github.com/async-profiler/async-profiler) downloaded locally

**Run on Mac**

```bash
java -agentpath:/Users/brian/files/async-profiler/async-profiler-4.2.1-macos/lib/libasyncProfiler.dylib=start,event=cpu,file=profile.html _1brc.java
```
