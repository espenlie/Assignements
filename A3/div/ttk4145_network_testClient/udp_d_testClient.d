import  std.stdio,
        std.conv,
        std.concurrency,
        std.algorithm,
        std.range,
        std.socket,
        core.thread;
        
        
        
void main(){
    writeln("version 02");
    
    string serverIP = "129.241.187.255";
    ushort serverPort = 20005;
    
    auto serverAddress = new InternetAddress(serverIP, serverPort);
        

    auto sendSock = new UdpSocket();
    sendSock.setOption(SocketOptionLevel.SOCKET, SocketOption.BROADCAST, 1);
    sendSock.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, 1);
    
    
    auto recvSock = new UdpSocket();
    recvSock.setOption(SocketOptionLevel.SOCKET, SocketOption.BROADCAST, 1);
    recvSock.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, 1);
    recvSock.bind(serverAddress);
    
    
    ubyte[1024] buf;
    
    while(true){
        Thread.sleep( 1.seconds );
        writeln("sending...");
        sendSock.sendTo("Hello UPD world", serverAddress);
        buf.clear;
        writeln(recvSock.receiveFrom(buf));
        writeln(cast(string)buf);
    } 
}

