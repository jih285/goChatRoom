package main

import (
    "fmt"
    "net"
    "strings"
    "time"
)

//var ConnMap map[string]*net.TCPConn
type client struct{
    name string
    ssocket *net.TCPConn
    currentRoom string
    myrooms []string
}
type chatroom struct{
    name string
    roomMember []*client
    log string
    lastActive time.Time
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
        fmt.Println("client disconnected")
        return -1
    }
    return 1
}

func doEveryWeek(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}


func checkExpiredRoom(tt time.Time){
	now:=time.Now()
	for _,room := range rooms{
		if timestample:=room.lastActive.Add(time.Hour * 24 * 7); timestample.Before(now){
			fmt.Println(room.name+" is expired! It has been deleted!")
			delete(rooms,room.name)
		}else{
			fmt.Println(tt)
		}
	}
}

func remove(s []string, i int) []string {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}
func removeClient(s []*client, i int) []*client {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

func say(myclient *client) {
    for {
        //fetch msg from a client
        data := make([]byte, 128)
        total, err := myclient.ssocket.Read(data)

        flag := checkErr(err)
        if flag <0 {
            break
        }

        fmt.Println(string(data[:total]), err)
        words := strings.Fields(string(data[:total]))
        iscommand:=false
        msg:=""
        if total<3 {
        	words=[]string{"please"}
        }
        switch command:=words[0]; command {
            case "jrename":
                myclient.name=words[1]
                msg="your name has been reset as: "+words[1]
                iscommand=true
            case "jcreate":
                var members []*client
                rooms[words[1]]=&chatroom{words[1],members,"",time.Now()}
                msg="room: "+words[1]+" has been created"
                iscommand=true
            case "jshowrooms":
                RoomList:="rooms: "
                for key,_ := range rooms{
                    RoomList+=key+" "
                }
                msg=RoomList;
                iscommand=true
            case "jshowmyrooms":
            	msg="your rooms are: "
            	if len(myclient.myrooms)>0{
	            	for _,myroomlist := range myclient.myrooms{
	            		msg+=myroomlist+" "
	            	}
            	}
            	iscommand=true
            case "jjoin":
                if _,ifexist:=rooms[words[1]]; ifexist {
                    rooms[words[1]].roomMember=append(rooms[words[1]].roomMember,myclient)
                    myclient.myrooms=append(myclient.myrooms,words[1])
                    msg="you have join room: "+words[1]+", now you could receive msg from this room\n------------------chat log------------------------\n"+rooms[words[1]].log
                }else{
                    msg="no such room, please check again"
                }
                iscommand=true

            case "jswitch":
               
                if  _, ok := rooms[words[1]]; ok {
                    for _,name := range myclient.myrooms{
                        if strings.Compare(name,words[1])==0{
                            myclient.currentRoom=words[1]
                            msg="you have switched to room: "+words[1]+", you could say something in this room now"
                        }else{
                            msg="you should join this room first"
                        }
                    }
                }else{
                    msg="no such room, please check again"
                }
                iscommand=true
            /*
            case "jcheck":{
            	iscommand=true
            	if _,ifexist:=rooms[words[1]]; ifexist {
                    
                    fmt.Println(rooms[words[1]].lastActive)
                }
            }
            case "jtest":{
            	checkExpiredRoom(time.Now());
            }
            */
            case "jleave":{
                iscommand=true
                findroom:=false
                for i,myroom:=range myclient.myrooms{
                    if strings.Compare(myroom,words[1])==0 {
                        myclient.myrooms=remove(myclient.myrooms,i)
                        findroom=true
                        for j,r := range rooms[words[1]].roomMember{
                            if r.ssocket.RemoteAddr().String() ==myclient.ssocket.RemoteAddr().String() {
                                rooms[words[1]].roomMember=removeClient(rooms[words[1]].roomMember,j)
                            }
                        }
                    if strings.Compare(myclient.currentRoom,words[1])==0{
                        myclient.currentRoom="none"
                    }
                    msg="you have left room: "+words[1]
                    break
                    }
                }
                if !findroom{
                    msg="you have to join this room, then leave"
                }
            }
            default:

            }


        bmsg:=[]byte(msg)

        //send system msg to the user itself
        if iscommand {
            myclient.ssocket.Write(bmsg)
        }else{
        	msg:=myclient.name+" from ["+myclient.currentRoom+"] said: "+string(data[:total])
        	 bmsg:=[]byte(msg)
            if strings.Compare(myclient.currentRoom, "none")!=0 {
                for _, conn := range rooms[myclient.currentRoom].roomMember {
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
                //record chat history
                rooms[myclient.currentRoom].log=rooms[myclient.currentRoom].log+msg+"\n"
                //update timestample of this chat room
                rooms[myclient.currentRoom].lastActive=time.Now()
            }else{
                msg="you should switch to a room first"
                bmsg:=[]byte(msg)
                myclient.ssocket.Write(bmsg)
            }

        }

    }
}

func main() {
    tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
    tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
    //ConnMap = make(map[string]*net.TCPConn)
    //var ConnMap make([]net.TCPConn, 11)
    rooms = make(map[string]*chatroom)
    //check every day if any room has no active member for 7 days
    go doEveryWeek(time.Hour * 24, checkExpiredRoom)
    for {
        tcpConn, _ := tcpListener.AcceptTCP()
        defer tcpConn.Close()
        //ConnMap=append(ConnMap,tcpConn)
        //ConnMap[tcpConn.RemoteAddr().String()] = tcpConn
        //tem:=[]string
        //welcomeMsg:="Hello welcome to Ji's chat lobby.\n users online now.\n\nYou could choose a chat room to join or create a new one.\n"+"useful commands:\n 1. jcreate [roomname]\n 2. jjoin [roomname] //receive messages from a chatroom, but can't speak\n 3. jswitch [roomname] //after join  a chatroom, you can siwtch to this room to say something\n 4. jleave [roomname]//you won't receive any message from this room\n 5. jshowrooms //show all existing rooms\n 6. jshowmyrooms //show the rooms that you join\n 7. rename [newname]";
        //tcpConn.Write([]byte(welcomeMsg))
        clients=append(clients,&client{tcpConn.RemoteAddr().String(),tcpConn,"none",[]string{}})
        fmt.Println("new client from:", tcpConn.RemoteAddr().String())
        go say(clients[len(clients)-1])
    }
}
