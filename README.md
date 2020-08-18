# syncmap 
![https://godoc.org/github.com/a8m/syncmap](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)
![LICENSE](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)
[![Build Status](https://travis-ci.com/a8m/syncmap.svg?token=ckAPcX3LvhP9wJPS6sgW&branch=master)](https://travis-ci.com/a8m/syncmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/a8m/syncmap)](https://goreportcard.com/report/github.com/a8m/syncmap)

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
  Or:
  ```bash
  $ go run github.com/a8m/syncmap -name IntMap "map[int]int"
  ```
  
2. Using `go generate`.
    
   - Add a directive with map definition:
     ```go
     //go:generate go run github.com/a8m/syncmap -name WriterMap map[string]io.Writer
   
     //go:generate go run github.com/a8m/syncmap -name Requests map[string]*http.Request
     ```
   - Then, run `go generate` on this package. 

   See [testdata/gen.go](https://github.com/a8m/syncmap/blob/master/testdata/gen.go) for more examples.
   
### How does it work?

`syncmap` didn't copy the code of `sync/map.go` and replace its identifiers. Instead, it reads the `sync/map.go` from
your `GOROOT`, parses it into an `*ast.File`, and runs a few mutators that bring it to the desired state.
Check the [code](https://github.com/a8m/syncmap/blob/master/syncmap.go#L91) for more information.

__How can we make sure it will continue to work?__ - I'm running a daily CI test on _TravisCI_.
   
### Benchmark
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

BenchmarkLoadOrStoreBalanced/*main.RWMutexMap-8     	 3000000	        400 ns/op
BenchmarkLoadOrStoreBalanced/*sync.Map-8            	 3000000	        400 ns/op
BenchmarkLoadOrStoreBalanced/*main.IntMap-8         	 5000000	        233 ns/op <--

BenchmarkLoadOrStoreUnique/*main.RWMutexMap-8       	 2000000	        744 ns/op
BenchmarkLoadOrStoreUnique/*sync.Map-8              	 2000000	        903 ns/op
BenchmarkLoadOrStoreUnique/*main.IntMap-8           	 3000000	        388 ns/op <--

BenchmarkLoadOrStoreCollision/*main.DeepCopyMap-8   	200000000	        7.29 ns/op
BenchmarkLoadOrStoreCollision/*main.RWMutexMap-8    	20000000	        97.5 ns/op
BenchmarkLoadOrStoreCollision/*sync.Map-8           	200000000	        9.11 ns/op
BenchmarkLoadOrStoreCollision/*main.IntMap-8        	500000000	        3.14 ns/op <--

BenchmarkRange/*main.DeepCopyMap-8                  	  500000	        4479 ns/op
BenchmarkRange/*main.RWMutexMap-8                   	   30000	        56834 ns/op
BenchmarkRange/*sync.Map-8                          	  300000	        4464 ns/op
BenchmarkRange/*main.IntMap-8                       	1000000000	        2.38 ns/op <--

BenchmarkAdversarialAlloc/*main.DeepCopyMap-8       	 2000000	        826 ns/op
BenchmarkAdversarialAlloc/*main.RWMutexMap-8        	20000000	        73.6 ns/op
BenchmarkAdversarialAlloc/*sync.Map-8               	 5000000	        303 ns/op
BenchmarkAdversarialAlloc/*main.IntMap-8            	10000000	        182 ns/op <--

BenchmarkAdversarialDelete/*main.DeepCopyMap-8      	10000000	        204 ns/op
BenchmarkAdversarialDelete/*main.RWMutexMap-8       	20000000	        78.3 ns/op
BenchmarkAdversarialDelete/*sync.Map-8              	20000000	        72.2 ns/op
BenchmarkAdversarialDelete/*main.IntMap-8           	100000000	        14.2 ns/op <--
```

Running benchmark with `-benchmem`
```
BenchmarkLoadMostlyHits/*main.DeepCopyMap-8         100000000	  12.7 ns/op	  7 B/op	  0 allocs/op
BenchmarkLoadMostlyHits/*main.RWMutexMap-8          30000000	  53.6 ns/op	  7 B/op	  0 allocs/op
BenchmarkLoadMostlyHits/*sync.Map-8                 100000000	  16.3 ns/op	  7 B/op	  0 allocs/op
BenchmarkLoadMostlyHits/*main.IntMap-8              200000000	  6.02 ns/op	  0 B/op	  0 allocs/op <--

BenchmarkLoadMostlyMisses/*main.DeepCopyMap-8       200000000	  7.99 ns/op	  7 B/op	  0 allocs/op
BenchmarkLoadMostlyMisses/*main.RWMutexMap-8        30000000	  52.6 ns/op	  7 B/op	  0 allocs/op
BenchmarkLoadMostlyMisses/*sync.Map-8               200000000	  8.87 ns/op	  7 B/op	  0 allocs/op
BenchmarkLoadMostlyMisses/*main.IntMap-8            1000000000	  2.88 ns/op	  0 B/op	  0 allocs/op <--

BenchmarkLoadOrStoreBalanced/*main.RWMutexMap-8     3000000	  357 ns/op	  71 B/op	  2 allocs/op
BenchmarkLoadOrStoreBalanced/*sync.Map-8            3000000	  417 ns/op	  70 B/op	  3 allocs/op
BenchmarkLoadOrStoreBalanced/*main.IntMap-8         5000000	  202 ns/op	  42 B/op	  1 allocs/op <--

BenchmarkLoadOrStoreUnique/*main.RWMutexMap-8       2000000	  648 ns/op	  178 B/op	  2 allocs/op
BenchmarkLoadOrStoreUnique/*sync.Map-8              2000000	  745 ns/op	  163 B/op	  4 allocs/op
BenchmarkLoadOrStoreUnique/*main.IntMap-8           3000000	  368 ns/op	  74 B/op	  2 allocs/op <--

BenchmarkLoadOrStoreCollision/*main.DeepCopyMap-8   300000000	  5.90 ns/op	  0 B/op	  0 allocs/op
BenchmarkLoadOrStoreCollision/*main.RWMutexMap-8    20000000	  94.5 ns/op	  0 B/op	  0 allocs/op
BenchmarkLoadOrStoreCollision/*sync.Map-8           200000000	  7.55 ns/op	  0 B/op	  0 allocs/op
BenchmarkLoadOrStoreCollision/*main.IntMap-8        1000000000	  2.68 ns/op	  0 B/op	  0 allocs/op <--

BenchmarkRange/*main.DeepCopyMap-8                  500000	  3376 ns/op	  0 B/op	  0 allocs/op
BenchmarkRange/*main.RWMutexMap-8                   30000	  56675 ns/op	  16384 B/op	  1 allocs/op
BenchmarkRange/*sync.Map-8                          500000	  3587 ns/op	  0 B/op	  0 allocs/op
BenchmarkRange/*main.IntMap-8                       2000000000	  1.75 ns/op	  0 B/op	  0 allocs/op <--

BenchmarkAdversarialAlloc/*main.DeepCopyMap-8       2000000	  761 ns/op	  535 B/op	  1 allocs/op
BenchmarkAdversarialAlloc/*main.RWMutexMap-8        20000000	  67.9 ns/op	  8 B/op	  1 allocs/op
BenchmarkAdversarialAlloc/*sync.Map-8               5000000	  264 ns/op	  51 B/op	  1 allocs/op
BenchmarkAdversarialAlloc/*main.IntMap-8            10000000	  176 ns/op	  28 B/op	  0 allocs/op <--

BenchmarkAdversarialDelete/*main.DeepCopyMap-8      10000000	  194 ns/op	  168 B/op	  1 allocs/op
BenchmarkAdversarialDelete/*main.RWMutexMap-8       20000000	  76.9 ns/op	  25 B/op	  1 allocs/op
BenchmarkAdversarialDelete/*sync.Map-8              20000000	  60.8 ns/op	  18 B/op	  1 allocs/op
BenchmarkAdversarialDelete/*main.IntMap-8           100000000	  13.1 ns/op	  0 B/op	  0 allocs/op <--
```


## LICENSE
I am providing code in the repository to you under MIT license. Because this is my personal repository, the license you receive to my code is from me and not my employer (Facebook)

