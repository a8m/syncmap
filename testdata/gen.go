package main

//go:generate syncmap -name Requests map[string]*http.Request

//go:generate syncmap -name StringMap map[string]interface{}

//go:generate syncmap -name WriterMap map[string]io.Writer

//go:generate syncmap -name stringerMap "map[string]interface{ String() string }"

//go:generate syncmap -name IntMap map[int]int

//go:generate syncmap -name StructMap "map[struct{ Name string }]struct{ Age int }"

//go:generate syncmap -name IntPtrs map[*int]*int

//go:generate syncmap -name StringByteChan "map[string]chan []byte"

//go:generate syncmap -name StringIntChan "map[string]chan int"

