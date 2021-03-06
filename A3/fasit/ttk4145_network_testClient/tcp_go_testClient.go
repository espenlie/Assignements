package main

import (
    ."fmt"
    ."net"
    ."strings"
    "time"
)

const (
    Single = iota
    Multi
)

type sendMessage_s struct {
    sendType    int
    msg         string
}

func receiver(conn Conn, receivedMsgs_c chan string){
    var buf [1024]byte
    for {
        _, err := conn.Read(buf[0:])
        if err != nil {
            Println(err)
            return
        }
        receivedMsgs_c <- Trim(string(buf[0:]), "\000")
    }
}

func messageGenerator(generatedMsgs_c chan sendMessage_s){
    generatedMsgs_c <- sendMessage_s{Single, "hello from go"}
    generatedMsgs_c <- sendMessage_s{Single, "Connect to: 129.241.187.159:22222"}
    time.Sleep(100*time.Millisecond)
    generatedMsgs_c <- sendMessage_s{Multi, "multi-test #1"}
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
//    initConn, err := Dial("tcp", "10.0.0.80:34933")
    initConn, err := Dial("tcp", "129.241.187.161:34933")
    if err != nil {
        Println(err)
    }
    listenAddr, _ := ResolveTCPAddr("tcp", ":22222")
    listenConn, _ := ListenTCP("tcp", listenAddr)

    connMap := make(map[string]Conn)

    receivedMsgs_c   := make(chan string)
    generatedMsgs_c  := make(chan sendMessage_s)
    newConn_c        := make(chan Conn, 10)

    newConn_c <- initConn
    go messageGenerator(generatedMsgs_c)
    go listener(listenConn, newConn_c)


    for {
        select {
        case msg := <- receivedMsgs_c:
            Println("   Main received: ", msg, " (end)")

        case sendMsg := <- generatedMsgs_c:
            if sendMsg.sendType == Multi {
                for k := range connMap {
                    connMap[k].Write(append([]byte(sendMsg.msg), []byte{0}...))
                }
            } else if sendMsg.sendType == Single {
                initConn.Write(append([]byte(sendMsg.msg), []byte{0}...))
            }

        case newConn := <- newConn_c:
            connMap[newConn.LocalAddr().String()] = newConn
            go receiver(newConn, receivedMsgs_c)
        }
    }
}















