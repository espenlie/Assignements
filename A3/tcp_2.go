package main

import (
    ."fmt"
    ."net"
    ."strings"
    "time"
)

func receiver(connection Conn, received chan string){
    var buf [1024]byte
    for {
        _, err := connection.Read(buf[0:])
        if err != nil {
            Println(err)
            return
        }
        received <- Trim(string(buf[0:]), "\000")
    }
}

func main(){
    init, err := Dial("tcp", "129.241.187.161:34933")
    if err != nil {
        Println(err)
    }
    messages := make(chan string)
    address, _ := ResolveTCPAddr("tcp", ":6969")
    listenConn, _ := ListenTCP("tcp", address)

    var sendMessage string

    go receiver(init, messages)
    Println(<-messages)
    connectMessage := "Connect to: 129.241.187.144:6969"
    init.Write(append([]byte(connectMessage), []byte{0}...))
    newConn, _ := listenConn.Accept()
    for {
        go receiver(newConn, messages)
        Println(<-messages)
        Scanf("%s", &sendMessage)
        newConn.Write(append([]byte(sendMessage), []byte{0}...))
        time.Sleep(100*time.Millisecond)
    }
}















