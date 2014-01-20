import  std.stdio,
        std.conv,
        std.range,
        std.algorithm,
        std.socket,
        std.concurrency,
        std.datetime,
        std.typecons,
        core.thread;

version = connect_to;
//version = simon_says;



void main(){


    string serverAddr = "129.241.187.161";
    //string serverAddr = "10.0.0.80";
    ushort port = 33546;

    auto sock = new TcpSocket(new InternetAddress(serverAddr, port));


    alias Tuple!(Socket, "socket", Tid, "receiver") conn_t;
    conn_t[] conns;

    conns ~= conn_t(    sock,
                        spawn( &Receive_Thr, thisTid, cast(shared)sock )   );
    auto messageGenerator   = spawn( &MessageGenerator_Thr );
    auto accepters          = [ spawn( &Accept_Thr, cast(ushort)22223 ),
                                spawn( &Accept_Thr, cast(ushort)22224 ),
                                spawn( &Accept_Thr, cast(ushort)22225 ),
                                spawn( &Accept_Thr, cast(ushort)22226 ),
                                spawn( &Accept_Thr, cast(ushort)22227 ) ];

    while(true){
        receive(
            (Tid t, send_msg sm){
                //if(t == messageGenerator){
                    if(sm.multi){
                        foreach(conn; conns){ conn.socket.send(sm.msg~"\0"); }
                    } else {
                        conns[0].socket.send(sm.msg~"\0");
                    }
                //}
            },
            (Tid t, receive_msg rm){
                if(conns.map!(a=>a.receiver).canFind(t)){
                    writeln("Main received: ", rm.msg);

                    if(rm.msg.skipOver("Simon says ")){
                        thisTid.send(thisTid, send_msg(false, rm.msg));
                    }

                }
            },
            (Tid t, shared Socket s){
                if(accepters.canFind(t)){
                    conns ~= conn_t(    cast(Socket)s,
                                        spawn( &Receive_Thr, thisTid, s )   );
                }
            },
            (Tid t, shutdown_msg s){
                conns.remove(conns.map!(a=>a.receiver).countUntil(t));
            },
            (Variant v){
                v.writeln(" (unknown)");
            }
        );
    }
}

void MessageGenerator_Thr(){
    version(connect_to){
        foreach(i; 0..5){
            Thread.sleep(100.msecs);
            ownerTid.send(thisTid, send_msg(false, "test# "~i.to!string));
        }
        immutable string localIP = new TcpSocket(new InternetAddress("www.google.com", 80)).localAddress.toAddrString;
        ownerTid.send(thisTid, send_msg(false, "Connect to: "~localIP~":22224"));
        ownerTid.send(thisTid, send_msg(false, "Connect to: "~localIP~":22225"));
        ownerTid.send(thisTid, send_msg(false, "Connect to: "~localIP~":22226"));
        ownerTid.send(thisTid, send_msg(false, "Connect to: "~localIP~":22227"));
        Thread.sleep(1.seconds);
        foreach(i; 1..4){
            ownerTid.send(thisTid, send_msg(true, "multi-test# "~i.to!string));
        }
    }
    version(simon_says){
        ownerTid.send(thisTid, send_msg(false, "Play Simon says"));
    }
}




void Receive_Thr(Tid recipient, shared Socket s){
    auto        sock    = cast(Socket)s;
    ubyte[1]    buff;
    string      buffer;

    writeln("    Receive thread started for ", sock.localAddress, " <-> ", sock.remoteAddress);

    while(sock.isAlive){
        while(sock.receive(buff) > 0){
            if(buff[0] != 0){
                buffer ~= buff;
            } else {
                recipient.send(thisTid, receive_msg(buffer));
                buffer.clear;
            }
        }

        if(sock.receive(buff) <= 0){
writeln("Disconnected... ", sock.localAddress, " <-> ", sock.remoteAddress);
            sock.shutdown(SocketShutdown.BOTH);
            sock.close();
            recipient.send(thisTid, shutdown_msg());
        }
    }
}


void Accept_Thr(ushort port){
    Socket acceptSock  = new TcpSocket();
    Socket newSock;

    acceptSock.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, 1);
    acceptSock.bind(new InternetAddress(port));
    acceptSock.listen(10);

    while(true){
        newSock = acceptSock.accept();
        writeln("    Accepted socket ", newSock.remoteAddress, " on ", Clock.currTime );
        ownerTid.send(thisTid, cast(shared)newSock);
    }


}


struct shutdown_msg {}
struct receive_msg {
    string msg;
}
struct send_msg {
    bool multi;
    string msg;
}

