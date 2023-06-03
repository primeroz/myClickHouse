# myClickHouse

All code tested on Linux

## Flow

```mermaid
flowchart LR
    A[DATAFILE] -->|Single Thread| B(READER)
    B -->|1 Row| C{WORKER1}
    B -->|1 Row| D{WORKER2}
    B -->|1 Row| E{WORKER3}
    B -->|1 Row| F{WORKER4}
    C -->G[Queue TOP10]
    D -->G[Queue TOP10]
    E -->G[Queue TOP10]
    F -->G[Queue TOP10]
```

## Test Data

There is already a test data file in `data/output.txt` of roughly 400K lines

To generate a different data set (with 1M elements) run
```
cd data
./generate-data.sh 1000000
```

## Parser

To run the parser with the `output.txt` file in the data directory
```
cd parser
make run
```

to Run the parser with another file as input 
```
cd parser
make build
echo "PATH_TO_THE_FILE" | ./parser
```

## Tests
To run the simple test
```
cd parser
make clean
make test
```
