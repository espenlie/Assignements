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
    conn := make(chan Conn, 10)
    conn <- init
    messages := make(chan string)
    conn2 := <- conn
    var sendMessage string
    for {
        go receiver(conn2, messages)
        Println(<-messages)
        Scanf("%s", &sendMessage)
        init.Write(append([]byte(sendMessage), []byte{0}...))
        time.Sleep(100*time.Millisecond)
    }
}

