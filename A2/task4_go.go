package main
 
import (
    . "fmt" // Using '.' to avoid prefixing functions with their package names
    . "runtime" // This is probably not a good idea for large projects...
)

var i = 0


func adder(intChannel chan int, finChannel chan string) {
    for x := 0; x < 1000000; x++ {
        j := <-intChannel
        i++
        intChannel <- j
    }
    finChannel <- "Add Done"
}
func subber(intChannel chan int, finChannel chan string) {
    for x := 0; x < 1000010; x++ {
        j := <-intChannel
        i--
        intChannel <- j
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
    Println(<-finChannel)
    Println(<-finChannel)
    Println(i);
}
