import java.io.BufferedReader;
import java.io.FileReader;

public class _1brc {
  public static void main(String[] args) {
    var start = System.nanoTime();
    System.out.println("Running calculations");

    var stop = System.nanoTime();
    var totalSeconds = (start - stop) / 1_000_000_000;
    var minutes = totalSeconds / 60;
    var seconds = totalSeconds % 60;
    System.out.println("Took %s to run");

    if (minutes > 0) {
      System.out.println(String.format("Elapsed time: %s minutes and %s seconds", minutes, seconds));
    } else {
      System.out.println(String.format("Elapsed time: %s seconds", seconds));
    }
  }

  private static void V1() {
    try (var reader = new BufferedReader(new FileReader("../1brc/measurements.txt"))) {
      String line;
      while ((line = reader.readLine()) != null) {
        var lineParts = line.split(";");
      }
    } catch (Exception ex) {
      System.err.println("something bad happened :(");
      System.err.println(ex);
    }
  }
}