package main
 
import (
    . "fmt" // Using '.' to avoid prefixing functions with their package names
    . "runtime" // This is probably not a good idea for large projects...
)

var i = 0


func adder(intChannel chan int, finChannel chan string) {
    for x := 0; x < 1000000; x++ {
        i := <-intChannel
        i++
        intChannel <- i
    }
    finChannel <- "Add Done"
}
func subber(intChannel chan int, finChannel chan string) {
    for x := 0; x < 1000010; x++ {
        i := <-intChannel
        i--
        intChannel <- i
    }
    finChannel <- "Sub Done"
}

func main() {
    GOMAXPROCS(NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
    intChannel := make(chan int, 1)
    intChannel <- i
    finChannel := make(chan string,2)
    go adder(intChannel, finChannel) // This spawns adder() as a goroutine
    go subber(intChannel, finChannel) // This spawns subber() as a goroutine
    s := <-finChannel
    Println("Done", s)
    p := <-finChannel
    Println("Done", p)
    i := <-intChannel
    Println("Done:", i);
}
