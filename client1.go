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
    fmt.Println("please enter your name")
    nickname, _ := reader.ReadString('\n')
    //fmt.Scanln(&nickname)

    fmt.Println("your name is :", nickname)

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

