# 1 Billion Row Challenge (1BRC) - Python

Attempting to do the challenge in python was kinda rough. With a starting time of 14 minutes I opted to start writing out multiple versions and start running multiple versions. On my Mac this did end up heating up the Mac which is probably what caused the slower times in the later versions. As such the times listed will only be for the times captured in PC

With the attempts I also opted not to use pypy as the JIT compiler as I'm mostly searching for what could cause slows downs or speed ups in python. This does mean that all the times could be potentially lower

Something else that made the challenge harder to tackle was trying to run the profiler and visualize what exactly is taking all the time for execution. Setting up py-spy didn't go well on Mac and cProfile didnt't provide much information, unsure if this is because of the profiler or because of python

Ended the optimization attempts early as it wasn't really intuitive debugging the slow downs through the profiler. py-spy might have been better to use but with the inability to setup on Mac I won't be able to use it

## Setup, Run and Profile

```bash
python3 -m venv .venv
source .venv/bin/activate
pip3 install -r requirements.txt
python3 main.py
snakeviz cpu.prof
```

## Versions

### V1

Super simple implementation to start with using dictionaries to track city values min, max, total and count, then at the end calculating the mean while building the output string. Starting with an extremely slow 11 minutes 59 seconds where there were major bottle necks with the usage of the split method splitting the string by semi-colon and the replace method removing the new line in the string. These two methods more than likely do not account for the entire slowness as they were recorded to take up ~3 minutes of the entire ~12 minutes run, however the profiler did not provide a detailed enough breakdown to investigate further

#### Timings

```
Average time 11minutes 59seconds
<method 'split' of 'str' objects> 1minute 50seconds
<method 'replace' of 'str' objects> 1minute 25seconds
```

### V2

Mostly the same as V1 but with some easy optimizations (at least I thought it might have been optimizing it), using a defaultdict for measurements and removing the usage of string.replace. Python provides the defaultdict interface which allows you to do something like `city_data = measurements[city]` and city_data will either have the data for a existing city or automatically create a new dictionary with default values you set before hand. I also used the "min" and "max" python built-ins instead of the previous if statements searching for the city min and max temperatures. The usage of the min and max methods were detected in the profiler, but it's hard to tell if it had a negative affect compared to the if statements or not. This update lead to a slower time of 13 minutes 5 seconds

#### Timings

```
Average time 13minutes 5seconds
<method 'split' of 'str' objects> minute 53seconds
<built-in method builtins.max> 1minute
<built-in method builtins.min> 1minute
```

#### V1 Snippet

```python
measurements: Dict[str, Values] = {}

# ...

if temperature_float < measurements[city]["min"]:
    measurements[city]["min"] = temperature_float
if temperature_float > measurements[city]["max"]:
    measurements[city]["max"] = temperature_float
```

#### V2 Snippet

```python
measurements: Dict[str, Values] = defaultdict(
    lambda: {"count": 0, "sum": 0.0, "min": 0.0, "max": 0.0}
)

# ...
city_data = measurements[city]
city_data["min"] = min(city_data["min"], temperature_float)
city_data["max"] = max(city_data["max"], temperature_float)
```

### V3

Identical to V2 but updates the final output to use a list to append the values and then do a string join with the list elements instead of using a string to track and append to it. This lead to a slightly slower time of 13 minutes 25 seconds, however it seems the split method was slower which was not updated, the difference between V2 and V3 may be negligible

#### Timings

```
Average time 13minutes 5seconds
<method 'split' of 'str' objects> 1minute 53seconds
<built-in method builtins.max> 1minute
<built-in method builtins.min> 1minute
```

#### V2 Snippet

```python
output = "{"
for city in sorted(measurements):
    measurement = measurements[city]
    mean = f"{(measurement["sum"] / measurement["count"]):.1f}"
    output += (
        f"{city}={measurements[city]["min"]}/{mean}/{measurements[city]["max"]}"
    )
output += "}"
print(output)
```

#### V3 Snippet

```python
output_parts = []
for city in sorted(measurements):
    measurement = measurements[city]
    mean = f"{(measurement["sum"] / measurement["count"]):.1f}"
    output_parts.append(
        f"{city}={measurements[city]["min"]}/{mean}/{measurements[city]["max"]}"
    )
print(f"{{{", ".join(output_parts)}}}")
```

### V4

Mostly the same as V3 but uses a dataclass instead of a dictionary for the city stats. This lead to a even larger slow down compared to V3, leading to a time of 16 minutes 10 seconds. Now that the dictionary update is being done in a differnt method instead of the V4 method used to track, the timing of the update (incrementing to count, determining min, determing max, and calculating sum) did come up in the profile. This did end up hiding the min and max timings however

#### Timings

```
Average time 16minutes 10seconds
main.py(update) 6minutes 47seconds
<method 'split' of 'str' objects> 2minutes 1second
```

#### V3 Snippet

```python
measurements: Dict[str, Values] = defaultdict(
    lambda: {"count": 0, "sum": 0.0, "min": 0.0, "max": 0.0}
)

# ...

city_data = measurements[city]
city_data["count"] += 1
city_data["min"] = min(city_data["min"], temperature_float)
city_data["max"] = max(city_data["max"], temperature_float)
city_data["sum"] += temperature_float
```

#### V4 Snippet

```python
@dataclass
class CityStats:
    count: int = 0
    max: float = 0.0
    min: float = 0.0
    sum: float = 0.0

    def update(self, temperature: float):
        self.count += 1
        self.max = max(self.max, temperature)
        self.min = min(self.min, temperature)
        self.sum += temperature

# ...

measurements: Dict[str, Values] = defaultdict(CityStats)

# ...

measurements[city].update(temperature_float)
```

### V5

Same as V4 but specifies to read the measurements file with read only and utf-8 encoding, this did seem to provide some speed up now at 15 minutes 55 seconds, speed up of ~15 seconds. This can be just a negligible speed up however similar to the speed up between V2 and V3

#### Timings

```
Average time 15minutes 55seconds
main.py(update) 6 minutes 41 seconds
<method 'split' of 'str' objects> 2minutes
```

#### V4 Snippet

```python
with open(measurments_file_path) as measurements_file:
```

#### V5 Snippet

```python
with open(measurments_file_path, "r", encoding="utf-8") as measurements_file:
```
