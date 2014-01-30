package main

import (
    ."fmt"
    ."net"
    ."strings"
    "time"
    "regexp"
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

func simonCheck(message string){
    match := regexp.MustCompile("Simon says")
    if match.MatchString(message) == true {
        Printf("Match")
    } else {
        Printf("No match")
    }
}


func main(){
    init, err := Dial("tcp", "129.241.187.161:34933")
    if err != nil {
        Println(err)
    }
    messages := make(chan string)
    startMessage := "Play Simon says\n"
    init.Write(append([]byte(startMessage), []byte{0}...))
    var sendMessage string
    for {
        go receiver(init, messages)
        simonCheck(<-messages)
        Println(<-messages)
        Scanf("%s", &sendMessage)
        init.Write(append([]byte(sendMessage), []byte{0}...))
        time.Sleep(100*time.Millisecond)
    }
}

