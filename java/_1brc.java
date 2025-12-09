import java.io.BufferedReader;
import java.io.FileReader;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;

public class _1brc {
  public static void main(String[] args) {
    var start = System.nanoTime();
    System.out.println("Running calculations");

    V1();

    var stop = System.nanoTime();
    var totalSeconds = (stop - start) / 1_000_000_000;
    var minutes = totalSeconds / 60;
    var seconds = totalSeconds % 60;

    if (minutes > 0) {
      System.out.println(String.format("Elapsed time: %s minutes and %s seconds", minutes, seconds));
    } else {
      System.out.println(String.format("Elapsed time: %s seconds", seconds));
    }
  }

  private static class Values {
    public int count = 0;
    public double max = 0;
    public double min = 0;
    public double sum = 0;

    public Values(double max, double min, double sum) {
      this.count = 1;
      this.max = max;
      this.min = min;
      this.sum = sum;
    }
  }

  /**
   * Quick and ditry initial attempt just looping through all the lines, splitting
   * by semi-colon, parsing the double and using a HashMap to track unique cities.
   * 
   * Mac Average time 2minutes 48seconds
   * java/io/BufferedReader.readLine 6,117 samples
   * java/lang/Double.parseDouble 3,758 samples
   * java/lang/String.split 4,602 samples
   * java/util/HashMap.get 2,220 samples
   */
  private static void V1() {
    var values = new HashMap<String, Values>();

    try (var reader = new BufferedReader(new FileReader("../1brc/measurements.txt"))) {
      String line;
      while ((line = reader.readLine()) != null) {
        var parts = line.split(";");
        var key = parts[0];
        var value = Double.parseDouble(parts[1]);

        var val = values.get(key);
        if (val == null) {
          values.put(key, new Values(value, value, value));
        } else {
          val.count++;
          val.sum += value;
          if (val.max < value)
            val.max = value;
          if (val.min > value)
            val.min = value;
        }
      }
    } catch (Exception ex) {
      System.out.println("Something went wrong :(");
      System.exit(1);
    }

    var keys = new ArrayList<String>(values.keySet());
    Collections.sort(keys);

    var output = "{";
    var idx = 0;
    var count = keys.size() - 1;
    for (var key : keys) {
      var value = values.get(key);
      output += String.format("%s=%.1f/%.1f/%.1f", key, value.min, value.sum / value.count, value.max);
      if (idx < count) {
        output += ", ";
      }
      idx++;
    }
    output += "}";
    System.out.println(output);
  }
}