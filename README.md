# syncmap

A __typed__ implementation of the Go `sync.Map` using code generation. 

### Install

```
go get -u github.com/a8m/syncmap
```

### Examples:

1. Using CLI
  ```bash
  $ syncmap -name IntMap "map[int]int"
  $ syncmap -name RequestMap -pkg mypkg "map[string]*http.Request"
  ```
  
2. Using `go generate`.
    
   - Add a directive with map definition:
     ```go
     //go:generate syncmap -name WriterMap map[string]io.Writer
   
     //go:generate syncmap -name Requests map[string]*http.Request
     ```
   - Then, run `go generate` on this package. 

   See [testdata/gen.go](https://github.com/a8m/syncmap/blob/master/testdata/gen.go) for more examples.
   
### Benchmarks
Benchmark tests were taken from the `sync` package.
```
BenchmarkLoadMostlyHits/*main.DeepCopyMap-8         	100000000	        15.1 ns/op
BenchmarkLoadMostlyHits/*main.RWMutexMap-8          	30000000	        54.4 ns/op
BenchmarkLoadMostlyHits/*sync.Map-8                 	100000000	        14.0 ns/op
BenchmarkLoadMostlyHits/*main.IntMap-8              	300000000	        5.65 ns/op <--

BenchmarkLoadMostlyMisses/*main.DeepCopyMap-8       	200000000	        10.2 ns/op
BenchmarkLoadMostlyMisses/*main.RWMutexMap-8        	30000000	        59.2 ns/op
BenchmarkLoadMostlyMisses/*sync.Map-8               	100000000	        11.3 ns/op
BenchmarkLoadMostlyMisses/*main.IntMap-8            	300000000	        4.05 ns/op <--

BenchmarkLoadOrStoreBalanced/*main.RWMutexMap-8     	 3000000	       400 ns/op
BenchmarkLoadOrStoreBalanced/*sync.Map-8            	 3000000	       400 ns/op
BenchmarkLoadOrStoreBalanced/*main.IntMap-8         	 5000000	       233 ns/op <--

BenchmarkLoadOrStoreUnique/*main.RWMutexMap-8       	 2000000	       744 ns/op
BenchmarkLoadOrStoreUnique/*sync.Map-8              	 2000000	       903 ns/op
BenchmarkLoadOrStoreUnique/*main.IntMap-8           	 3000000	       388 ns/op <--

BenchmarkLoadOrStoreCollision/*main.DeepCopyMap-8   	200000000	         7.29 ns/op
BenchmarkLoadOrStoreCollision/*main.RWMutexMap-8    	20000000	         97.5 ns/op
BenchmarkLoadOrStoreCollision/*sync.Map-8           	200000000	         9.11 ns/op
BenchmarkLoadOrStoreCollision/*main.IntMap-8        	500000000	         3.14 ns/op <--

BenchmarkRange/*main.DeepCopyMap-8                  	  500000	      4479 ns/op
BenchmarkRange/*main.RWMutexMap-8                   	   30000	     56834 ns/op
BenchmarkRange/*sync.Map-8                          	  300000	      4464 ns/op
BenchmarkRange/*main.IntMap-8                       	1000000000	      2.38 ns/op <--

BenchmarkAdversarialAlloc/*main.DeepCopyMap-8       	 2000000	       826 ns/op
BenchmarkAdversarialAlloc/*main.RWMutexMap-8        	20000000	      73.6 ns/op
BenchmarkAdversarialAlloc/*sync.Map-8               	 5000000	       303 ns/op
BenchmarkAdversarialAlloc/*main.IntMap-8            	10000000	       182 ns/op <--

BenchmarkAdversarialDelete/*main.DeepCopyMap-8      	10000000	         204 ns/op
BenchmarkAdversarialDelete/*main.RWMutexMap-8       	20000000	        78.3 ns/op
BenchmarkAdversarialDelete/*sync.Map-8              	20000000	        72.2 ns/op
BenchmarkAdversarialDelete/*main.IntMap-8           	100000000	        14.2 ns/op <--
```
