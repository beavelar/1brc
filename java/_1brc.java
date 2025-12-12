import java.io.BufferedReader;
import java.io.FileReader;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;

public class _1brc {
  public static void main(String[] args) {
    var start = System.nanoTime();
    System.out.println("Running calculations");

    // V1();
    // V2();
    V3();

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
   * Average time 1minute 34seconds
   * java/io/BufferedReader.readLine 3,545 samples
   * java/lang/Double.parseDouble 1,402 samples
   * java/lang/String.split 3,329 samples
   * java/util/HashMap.get 1,325 samples
   * 
   * Mac Average time 2minutes 48seconds
   * java/io/BufferedReader.readLine 6,117 samples
   * java/lang/Double.parseDouble 3,758 samples
   * java/lang/String.split 4,602 samples
   * java/util/HashMap.get 2,220 samples
   */
  public static void V1() {
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
      ex.printStackTrace();
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

  /**
   * Identical to V2 but increases the buffer size of BufferedReader to 12MB
   * instead of the default 8KB
   * 
   * Average time 1minute 34seconds
   * java/io/BufferedReader.readLine 3,740 samples
   * java/lang/Double.parseDouble 1,226 samples
   * java/lang/String.split 3,577 samples
   * java/util/HashMap.get 1,314 samples
   * 
   * Mac Average time 2minutes 43seconds
   * java/io/BufferedReader.readLine 5,271 samples
   * java/lang/Double.parseDouble 4,005 samples
   * java/lang/String.split 4,750 samples
   * java/util/HashMap.get 2,113 samples
   */
  public static void V2() {
    var values = new HashMap<String, Values>();

    // 12MB buffer
    try (var reader = new BufferedReader(new FileReader("../1brc/measurements.txt"), 12 * 1024 * 1024)) {
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
      ex.printStackTrace();
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

  /**
   * Updates the file reading portion to use FileReader and provide the contents
   * to a char array, iterate over that char array instead of reading through the
   * buffered reader
   * 
   * Average time 1minute 32seconds
   * java/io/Reader.read 1,847 samples
   * java/lang/Double.parseDouble 1,556 samples
   * java/lang/StringBuilder.append 1,354 samples
   * java/lang/StringBuilder.toString 1,649 samples
   * java/util/HashMap.get 1,648 samples
   */
  public static void V3() {
    var values = new HashMap<String, Values>();
    var cbufSize = 12 * 1024 * 1024;
    var cbuf = new char[cbufSize];

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

            var val = values.get(city);
            if (val == null) {
              values.put(city, new Values(temp, temp, temp));
            } else {
              val.count++;
              val.sum += temp;
              if (val.max < temp)
                val.max = temp;
              if (val.min > temp)
                val.min = temp;
            }
            continue;
          }

          if (parsingCity) {
            citySb.append(cbuf[idx]);
          } else {
            tempSb.append(cbuf[idx]);
          }
        }
      }
    } catch (Exception ex) {
      System.out.println("Something went wrong :(");
      ex.printStackTrace();
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