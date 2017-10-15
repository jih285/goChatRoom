package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
)

var ch chan int = make(chan int)

var nickname string

func reader(conn *net.TCPConn) {
    buff := make([]byte, 128)
    for {
        j, err := conn.Read(buff)
        if err != nil {
            ch <- 1
            break
        }

        fmt.Println(string(buff[0:j]))
    }
}

func main() {

    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
    conn, err := net.DialTCP("tcp", nil, tcpAddr)

    if err != nil {
        fmt.Println("Server is not starting")
        os.Exit(0)
    }

    //为什么不能放到if之前？ err不为nil的话就是painc了 (painc 与 defer 辨析一下！！！)
    defer conn.Close()
    go reader(conn)
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("Hello welcome to Ji's chat lobby.\nYou could choose a chat room to join or create a new one.\n"+"useful commands:\n 1. jcreate [roomname]\n 2. jjoin [roomname] //receive messages from a chatroom, but can't speak\n 3. jswitch [roomname] //after join  a chatroom, you can siwtch to this room to say something\n 4. jleave [roomname]//you won't receive any message from this room\n 5. jshowrooms //show all existing rooms\n 6. jshowmyrooms //show the rooms that you join\n 7. rename [newname]");
    fmt.Println("please enter your name")
    nickname, _ := reader.ReadString('\n')
    //fmt.Scanln(&nickname)
    fmt.Println("your name is :", nickname)

    setNameCMD:="jrename "+nickname
    b := []byte(setNameCMD)
    conn.Write(b)

    for {
        //var msg string
        //fmt.Scanln(&msg)
        msg, _ := reader.ReadString('\n')
        //b := []byte("<" + nickname + ">" + "said: " + msg)
        b := []byte(msg)

        conn.Write(b)

        //select 为非阻塞的
        select {
        case <-ch:
            fmt.Println("Server error, please reconnect!")
            os.Exit(1)
        default:
            //不加default的话，那么 <-ch 会阻塞for， 下一个输入就没有法进行
        }

    }
}

