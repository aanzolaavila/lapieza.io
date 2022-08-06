# lapieza.io Go technical test

## Build
To build this project you need GNU Make and go 1.18+ installed, after that do:
```bash
$ make build
```
That will yield two executables in the `bin` folder: `run` and `random`.

## Executables

### `run`
This is the main program to process prices and yield the maximum amount of earnings per stock window.

The input format is expected to be as an EOF terminanted file, the file can be specified with the `--file` option, if none is especified the default is standard input `stdin`.

Each line represents the stock window to be processed, each number in the line is separated by a space.

#### Sample input
```
7 1 5 3 6 4
7 6 4 3 1

```

#### How it works
This program reads the file line by line, each is converted into a slice of `int64` numbers, each of those is fed into a channel `inCh`.

If there is an error in reading the line of numbers, there is an `error` fed into the `errs` channel.

Concurrently running in other goroutine, we read from that channel `inCh`, calculate the answer with a linear `O(n)` algorithm, then feed it into another channel `outCh`.

Meanwhile in the main thread, we are reading from the `outCh` and `errs` channel and print them in a FIFO manner. This means that the is NO guarantee for the order of the output as some slices could be larger than others and take longer to process.

#### Running
```bash
$ ./run --file in.txt
$ cat in.txt | ./run
```
**The `in.txt` file is included in this repo and includes two (2) example inputs.**

### `random`
This program generated random numbers in the format specified by the main executable.

Its arguments are:
- `--min`
  Minimum value that the values in each line can have.
- `--max`
  Maximum value that the values in each line can have.
- `--arrays`
  Number of lines to generate.
- `--size`
  The amount of numbers in each line to generate.

#### Running
```bash
$ ./random --max 100 --arrays 1000000 --size 100 | ./run
```
