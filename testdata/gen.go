package main

//go:generate go run github.com/a8m/syncmap -name Requests map[string]*http.Request

//go:generate go run github.com/a8m/syncmap -name StringMap map[string]interface{}

//go:generate go run github.com/a8m/syncmap -name WriterMap map[string]io.Writer

//go:generate go run github.com/a8m/syncmap -name stringerMap "map[string]interface{ String() string }"

//go:generate go run github.com/a8m/syncmap -name IntMap map[int]int

//go:generate go run github.com/a8m/syncmap -name StructMap "map[struct{ Name string }]struct{ Age int }"

//go:generate go run github.com/a8m/syncmap -name IntPtrs map[*int]*int

//go:generate go run github.com/a8m/syncmap -name StringByteChan "map[string](chan []byte)"

//go:generate go run github.com/a8m/syncmap -name StringIntChan "map[string](chan int)"
