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