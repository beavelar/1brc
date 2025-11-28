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
# Mac Average 13minutes 31seconds
def V1():
  measurements: Dict[str, Values] = {}
  measurments_file_path = os.path.realpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, "1brc", "measurements.txt"))
  with open(measurments_file_path) as measurements_file:
    for line in measurements_file:
      [city, temperature] = line.split(";")
      temperature_float = float(temperature)
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

if __name__ == "__main__":
  profiler = cProfile.Profile()
  print("Running calculations")
  profiler.enable()
  start = time.perf_counter()
  
  V1()
  
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