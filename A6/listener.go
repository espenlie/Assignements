package main

import (
    ."fmt"
    ."net"
    ."strings"
    "os/exec"
    "time"
    "strconv"
)

func receiver(connection Conn, received chan string,status chan string){
    var buf [1024]byte
    for {
        _, err := connection.Read(buf[0:])
        if err != nil {
            status <- "YO"
            //Println(err)
            return
        }
        received <- Split(string(buf[0:]), "EOL")[0]
    }
}


func listener(conn *TCPListener, newConn_c chan Conn){
    for {
        newConn, err := conn.Accept()
        if err != nil {
            //Println(err)
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
    status_c        := make(chan string)

    var counter int

    go listener(listenConn, newConn_c)

    for {
        select {
        case <-status_c:
            listenConn.Close()
            cmd := exec.Command("mate-terminal","-x", "go", "run", "/home/student/LL/Assignements/A6/listener.go")
            err2 := cmd.Run()
            if err2 != nil {
                //Println(err2)
            }
            time.Sleep(1000 * time.Millisecond)
            initConn, err := Dial("tcp", "localhost:6969")
            if err != nil {
                //Println(err)
            }
            var sendMessage string
            for {
                    counter=counter+1
                    sendMessage=strconv.Itoa(counter)+"EOL"
                    //Println(strconv.Itoa(counter))
                    initConn.Write(append([]byte(sendMessage), []byte{0}...))
                    time.Sleep(1000 * time.Millisecond)
            }
        case msg := <-receivedMsgs_c:
            Println(msg)
            counter, _ = strconv.Atoi(msg)
        case newConn := <-newConn_c:
            connMap[newConn.LocalAddr().String()] = newConn
            go receiver(newConn, receivedMsgs_c, status_c)
        }
    }
}

