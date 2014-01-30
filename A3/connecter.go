package main

import (
    ."fmt"
    ."net"
)

const (
    Single = iota
    Multi
)

type sendMessage_s struct {
    sendType    int
    msg         string
}

func main(){
    initConn, err := Dial("tcp", "129.241.187.144:6969")
    if err != nil {
        Println(err)
    }

    var sendMessage string
    for {   
        Scanf("%s", &sendMessage)
        sendMessage=sendMessage+"EOL"
        initConn.Write(append([]byte(sendMessage), []byte{0}...))
    }
}
