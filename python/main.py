from collections import defaultdict
from dataclasses import dataclass
from typing import Dict, TypedDict
import cProfile
import os
import time

class Values(TypedDict):
    city: str
    count: int
    max: float
    min: float
    sum: str

# Super simple implementation to start with, just using dictionaries to track city values
# min, max, total and count, then at the end calculating the mean while building the output
# string
#
def V1():
  measurements: Dict[str, Values] = {}
  measurments_file_path = os.path.realpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, "1brc", "measurements.txt"))
  with open(measurments_file_path) as measurements_file:
    for line in measurements_file:
      [city, temperature] = line.split(";")
      temperature_float = float(temperature.replace("\n", ""))
      if city in measurements:
        measurements[city]["count"] += 1
        measurements[city]["sum"] += temperature_float

        if temperature_float < measurements[city]["min"]:
          measurements[city]["min"] = temperature_float
        if temperature_float > measurements[city]["max"]:
          measurements[city]["max"] = temperature_float
      else:
        measurements[city] = {
          "city": city,
          "count": 1,
          "max": temperature_float,
          "min": temperature_float,
          "sum": temperature_float
        }
    
    output = "{"
    for city in sorted(measurements):
      measurement = measurements[city]
      mean = f"{(measurement["sum"] / measurement["count"]):.1f}"
      output += f"{city}={measurements[city]["min"]}/{mean}/{measurements[city]["max"]}"
    output += "}"
    print(output)

# Mostly the same as V1 but with some easy optimizatinos, using a defaultdict for measurements
# and removing the usage of string.replace. Using defaultdict, we just need to access the key
# (ex. measurment[city]) and if the key exist, it'll return you the value, otherwise it'll
# return the default value listed
#
def V2():
  measurements: Dict[str, Values] = defaultdict(lambda: {"count": 0, "sum": 0.0, "min": 0.0, "max": 0.0})
  measurments_file_path = os.path.realpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, "1brc", "measurements.txt"))
  with open(measurments_file_path) as measurements_file:
    for line in measurements_file:
      [city, temperature] = line.split(";")
      temperature_float = float(temperature)

      city_data = measurements[city]
      city_data["count"] += 1
      city_data["min"] = min(city_data["min"], temperature_float)
      city_data["max"] = max(city_data["max"], temperature_float)
      city_data["sum"] += temperature_float
    
    output = "{"
    for city in sorted(measurements):
      measurement = measurements[city]
      mean = f"{(measurement["sum"] / measurement["count"]):.1f}"
      output += f"{city}={measurements[city]["min"]}/{mean}/{measurements[city]["max"]}"
    output += "}"
    print(output)

# Identical to V2 but updates the output string building in the end to append to a list and
# do a join with the values to avoid string concactinations
#
def V3():
  measurements: Dict[str, Values] = defaultdict(lambda: {"count": 0, "sum": 0.0, "min": 0.0, "max": 0.0})
  measurments_file_path = os.path.realpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, "1brc", "measurements.txt"))
  with open(measurments_file_path) as measurements_file:
    for line in measurements_file:
      [city, temperature] = line.split(";")
      temperature_float = float(temperature)

      city_data = measurements[city]
      city_data["count"] += 1
      city_data["min"] = min(city_data["min"], temperature_float)
      city_data["max"] = max(city_data["max"], temperature_float)
      city_data["sum"] += temperature_float
    
    output_parts = []
    for city in sorted(measurements):
      measurement = measurements[city]
      mean = f"{(measurement["sum"] / measurement["count"]):.1f}"
      output_parts.append(f"{city}={measurements[city]["min"]}/{mean}/{measurements[city]["max"]}")
    print(f"{{{", ".join(output_parts)}}}")

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

# Mostly the same as V3 but uses a dataclass instead of a dictionary for the city stats
def V4():
  measurements: Dict[str, Values] = defaultdict(CityStats)
  measurments_file_path = os.path.realpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, "1brc", "measurements.txt"))
  with open(measurments_file_path) as measurements_file:
    for line in measurements_file:
      [city, temperature] = line.split(";")
      temperature_float = float(temperature)
      measurements[city].update(temperature_float)
    
    output_parts = []
    for city in sorted(measurements):
      measurement = measurements[city]
      mean = f"{(measurement["sum"] / measurement["count"]):.1f}"
      output_parts.append(f"{city}={measurements[city]["min"]}/{mean}/{measurements[city]["max"]}")
    print(f"{{{", ".join(output_parts)}}}")

# Same as V4 but specifies to update the measurements file with read and utf-8 encoding
def V5():
  measurements: Dict[str, Values] = defaultdict(CityStats)
  measurments_file_path = os.path.realpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, "1brc", "measurements.txt"))
  with open(measurments_file_path, 'r', encoding='utf-8') as measurements_file:
    for line in measurements_file:
      [city, temperature] = line.split(";")
      temperature_float = float(temperature)
      measurements[city].update(temperature_float)
    
    output_parts = []
    for city in sorted(measurements):
      measurement = measurements[city]
      mean = f"{(measurement["sum"] / measurement["count"]):.1f}"
      output_parts.append(f"{city}={measurements[city]["min"]}/{mean}/{measurements[city]["max"]}")
    print(f"{{{", ".join(output_parts)}}}")

if __name__ == "__main__":
  profiler = cProfile.Profile()
  print("Running calculations")
  profiler.enable()
  start = time.perf_counter()
  
  # V1()
  # V2()
  # V3()
  # V4()
  V5()

  end = time.perf_counter()
  profiler.disable()
  profiler.dump_stats("cpu.prof")

  elapsed = end - start
  minutes = int(elapsed // 60)
  seconds = elapsed % 60

  if minutes > 0:
    print(f"Elapsed time: {minutes} minutes and {seconds:.2f} seconds")
  else:
    print(f"Elapsed time: {seconds:.2f} seconds")