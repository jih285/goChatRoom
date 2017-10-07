package main

import (
    "fmt"
    "net"
    "strings"
)

//var ConnMap map[string]*net.TCPConn
type client struct{
    name string
    ssocket *net.TCPConn
}
type chatroom struct{
    name string
    roomMember []*client
}

//var ConnMap []*net.TCPConn
var clients []*client   //store all of the clients
var rooms map[string]*chatroom


func checkErr(err error) int {
    if err != nil {
        if err.Error() == "EOF" {
            fmt.Println("client quit")
            return 0
        }
        fmt.Println("error")
        return -1
    }
    return 1
}

func say(myclient *client) {
    for {
        //fetch msg from a client
        data := make([]byte, 128)
        total, err := myclient.ssocket.Read(data)

        msg:=myclient.name+" said: "+string(data[:total])
        

        fmt.Println(string(data[:total]), err)
        words := strings.Fields(string(data[:total]))

        switch command:=words[0]; command {
            case "jrename":
                myclient.name=words[1]
            case "jcreate":
                var members []*client
                rooms[words[1]]=&chatroom{words[1],members}
                fmt.Println("room: "+words[1]+" has been created")
            case "jshowrooms":
                RoomList:="rooms: "
                for key,_ := range rooms{
                    RoomList+=key+" "
                }
                msg=RoomList;

            default:
                // freebsd, openbsd,
                // plan9, windows...
            }
        flag := checkErr(err)
        if flag == 0 {
            break
        }

        bmsg:=[]byte(msg)
        //broadcast 
        for _, conn := range clients {
            /*
            if conn.RemoteAddr().String() == myclient.ssocket.RemoteAddr().String() {
                //don't send msg back to its sender
                continue
            }
            */
            //msg:=[]byte(conn.name+" said: ")
           // msg=append(msg,data[:total])
            conn.ssocket.Write(bmsg)
        }
    }
}

func main() {
    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
    tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
    /*
        map 定义完后，还要make? (哪些数据类型定义完后，还要make?)
        http://stackoverflow.com/questions/27267900/runtime-error-assignment-to-entry-in-nil-map
    */
    //ConnMap = make(map[string]*net.TCPConn)
    //var ConnMap make([]net.TCPConn, 11)
    rooms = make(map[string]*chatroom)
    for {

        tcpConn, _ := tcpListener.AcceptTCP()
        defer tcpConn.Close()
        //ConnMap=append(ConnMap,tcpConn)
        //ConnMap[tcpConn.RemoteAddr().String()] = tcpConn
        clients=append(clients,&client{tcpConn.RemoteAddr().String(),tcpConn})
        fmt.Println("new client from:", tcpConn.RemoteAddr().String())

        go say(clients[len(clients)-1])
    }
}
