package main

import (
    ."fmt"
    ."net"
    ."strings"
)

func receiver(connection Conn, received chan string){
    var buf [1024]byte
    for {
        _, err := connection.Read(buf[0:])
        if err != nil {
            Println(err)
            return
        }
        received <- Split(string(buf[0:]), "EOL")[0]
    }
}


func listener(conn *TCPListener, newConn_c chan Conn){
    for {
        newConn, err := conn.Accept()
        if err != nil {
            Println(err)
        }
        newConn_c <- newConn
    }
}


func main(){

    listenAddr, _ := ResolveTCPAddr("tcp", ":6969")
    listenConn, _ := ListenTCP("tcp", listenAddr)

    connMap         := make(map[string]Conn)

    receivedMsgs_c  := make(chan string)
    newConn_c       := make(chan Conn, 10)

    go listener(listenConn, newConn_c)

    for {
        select {
        case msg := <- receivedMsgs_c:
            Println("   Main received: ", msg, " (end)")

        case newConn := <- newConn_c:
            connMap[newConn.LocalAddr().String()] = newConn
            go receiver(newConn, receivedMsgs_c)
        }
    }
}

