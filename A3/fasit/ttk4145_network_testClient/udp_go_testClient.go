package main;

import (
    ."fmt"
    ."net"
    "time"
)

var localIP = "129.241.187.156"

func main(){
    serverAddr, err := ResolveUDPAddr("udp4", "129.241.187.255:20005")
    if err != nil { Println(err) }
    conn, err := DialUDP("udp4", nil, serverAddr)
    if err != nil { Println(err) }
    
    go listener()
    
    for {
        conn.Write([]byte("Hello UDP go world"))
        time.Sleep(1*time.Second)
    }
    
//    Println(conn)
//    Println("done")
    
}


func listener(){

    serverAddr, err := ResolveUDPAddr("udp4", ":20005")
    if err != nil { Println(err) }
    conn, err := ListenUDP("udp4", serverAddr)
    if err != nil { Println(err) }
    
    var buf [512]byte
    
    for {
        _, remoteAddr, err := conn.ReadFromUDP(buf[0:])
        if err != nil { Println(err) }
        if localIP != remoteAddr.IP.String() {
            Println("Received:", string(buf[0:]))
        }
    }
}


