# Go profiling example

*Sample code I use for Go performance profiling talk. Code is unoptimized on purpose.*

Here we have a sample go application that creates a "blockchain" in memory, by generating random blocks and hashing their data together with previous block's hash.
The domain of the application is not important, what is important is that it involves a lot of allocations and computations that we can inspect using go `pprof` tool.

Profiling starts by importing `"runtime/pprof"` package (or `"net/http/pprof"` if you want to expose profiling endpoints via HTTP).

You can start CPU profiling at any time, and profile will be written to the file until tracing is stopped. Memory profile works buy saving momentarily snapshots.
In both cases you need to take care to create a file to store the trace, and close it after tracing is finished.

## CPU profiling

Create a file:

```
f, err := os.Create("pprof.out")
if err != nil {
    panic(err)
}
```

Start CPU profiling:

```
if err := pprof.StartCPUProfile(f); err != nil {
    panic(err)
}
```

Do some work and then stop the profiler and close the file we opened earlier:

```
pprof.StopCPUProfile()
f.Close()
```

Run the code:

```
$ go run main.go
```

This will create `pprof.out` file that we can examine using go `pprof` tool:

```
$ go tool pprof -http 127.0.0.1:8080 pprof.out
```

This will open web browser with graph of function calls with relative and absolute CPU time these functions consume.

## Memory profile

To profile the memory allocation we create a snapshot of current memory state in a certain period of time:

```
f2, err := os.Create("mem.out")
if err != nil {
    panic(err)
}

pprof.WriteHeapProfile(f2)
f2.Close()
```

Run the code:

```
$ go run main.go
```

Then run pprof tool the same way:

```
$ go tool pprof -http 127.0.0.1:8080 mem.out
```

## Tracing

Another useful tool is `trace`. It allows to look at your code execution timeline with nanosecond resolution. It is very useful to investigate allocation and garbage collector related issues.

One way to enable trace is to import `"runtime/trace"` and enable trace manually in the code. Also `"net/http/trace"` has a dedicated `/debug/pprof/trace` endpoint. Last and not least you can pass `-trace` flag to `go test` command. Let's do just excactly that:

```
$ go test -trace trace.out -bench . .
```

Now our test (in our case it is a benchmark) has generated a trace file, we can analyze it:

```
$ go tool trace trace.out
```

This will launch the browser where you can take a look at your trace.
