package main


import (
    . "fmt" // Using '.' to avoid prefixing functions with their package names
    . "runtime" // This is probably not a good idea for large projects...
    . "time"
    "sync"
)

type Counter struct {
    mu sync.Mutex
    i int64
}

//var i = 0

func (c *Counter) adder() {
    for x := 0; x < 1000000; x++ {
        c.mu.Lock()
        c.i++
        c.mu.Unlock()
    }
}
func (c *Counter) subber() {
    for x := 0; x < 1000000; x++ {
        c.mu.Lock()
        c.i--
        c.mu.Unlock()
    }
}
//func subber() {
//    for x := 0; x < 1000000; x++ {
//        i--
//    }
//}

func main() {
    GOMAXPROCS(NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
    var views Counter
    go adder() // This spawns adder() as a goroutine
    go subber() // This spawns subber() as a goroutine
    Sleep(100*Millisecond)
    Println("Done:", i);
}
