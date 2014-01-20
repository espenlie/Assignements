import  std.algorithm,
        std.range,
        std.conv,
        std.stdio,
        std.socket,
        std.string,
        std.concurrency,
        std.process,
        core.thread,
        core.sys.posix.unistd,
        std.c.stdlib;




immutable ushort[] ports = iota(20000, 20021, 1).map!(to!ushort).array;


int main(string[] args){
    string programName = args[0][2..$];
    int[] otherPids;
    
    
    writeln("Looking for other running servers...");
    writeln("ThisPID: ", thisProcessID);

    do {
        version(linux){
            otherPids = executeShell("pgrep \"" ~ programName ~ "\"")
                        .output
                        .split
                        .to!(int[])
                        .sort
                        .setDifference([thisProcessID])
                        .array;
            assert(otherPids.canFind(thisProcessID) == false);
            Thread.sleep(5.seconds);
        } else {
            pragma(msg, "\n    Automatic restart not supported, compiling without\n");
        }
    } while(otherPids.length);
    
    
    writeln("This is now the active server");
    
    //////////////////////////////////////
    
    immutable string localIP = new TcpSocket(new InternetAddress("www.google.com", 80)).localAddress.toAddrString;
    immutable string bcastIP = localIP[0 .. localIP.lastIndexOf(".")+1] ~ "255";
    writeln("\n    Local IP is: ", localIP, "\n");
    

	foreach(port; ports){
	    spawnLinked(&UDPEcho_task, localIP, bcastIP, cast(ushort)port);
	}

    receive(
        (LinkTerminated lt){
            writeln("A thread has died!");
            spawnProcess("./" ~ programName);
            exit(17);
        }
    );

    return 0;
}



void UDPEcho_task(string localIP, string bcastIP, ushort port){
    writeln("Echo_task for port ", port, " started");
    
    auto serverAddress = new InternetAddress(bcastIP, port);
        

    auto sendSock = new UdpSocket();
    sendSock.setOption(SocketOptionLevel.SOCKET, SocketOption.BROADCAST, 1);
    sendSock.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, 1);
    
    
    auto recvSock = new UdpSocket();
    recvSock.setOption(SocketOptionLevel.SOCKET, SocketOption.BROADCAST, 1);
    recvSock.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, 1);
    recvSock.bind(serverAddress);
    
    
    ubyte[1024] buf;
    Address remoteAddr;
    
    
    while(recvSock.receiveFrom(buf, remoteAddr) > 0){
        if(remoteAddr.toAddrString != localIP){
            writeln(port, " received from ", remoteAddr, " : ", (cast(string)buf).strip('\0'));
            sendSock.sendTo(cast(ubyte[])"You said: " ~ buf, serverAddress);
        }
        buf.clear;
    }
}












