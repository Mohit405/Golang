package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var TodoList []string
func cmd(message string)string{
    inputArr:=strings.Split(message," ")
    return inputArr[0];
}
func convert(message string)string{
    inputArr:=strings.Split(message," ")
    var ans string
    for i:=1;i<len(inputArr);i++{
        ans+=inputArr[i]
    }
    return ans
}
func UpdateTodoList(value string){
    for i,num:=range TodoList{
        if num==value{
            TodoList = append(TodoList[:i],TodoList[i+1:]... )
            break
        }
    }
}
func main(){

    //setting up a http server to handle the client request 
    http.HandleFunc("/todo",func(w http.ResponseWriter,r* http.Request){
        //converting this request into the websocket protocol to make it a full-duplex
        //it can be done using the Upgrader function of the websocket module
        conn,err:=upgrader.Upgrade(w,r,nil)
        if err!=nil{
            log.Print("Websocket error",err)
            return
        }
        for{
            //reading from the server
            mtype,message,err:=conn.ReadMessage()
            if err!=nil{
                log.Println("Connection error",err)
                return
            }
            input:=string(message)
            command:=cmd(input)
            msg:=convert(input)
            if command=="add"{
                TodoList = append(TodoList, msg)
            }else if command=="done"{
                UpdateTodoList(msg)
            }
            output:="Todos:\n"
            for _,value:=range TodoList{
                output+="\n"+value+"\n"
            }
            output+="\n--------------------------------"
            message=[]byte(output)
            err=conn.WriteMessage(mtype,message)
            if err!=nil{
                log.Println("Connection err",err)
                return
            }
        }
    })
    http.HandleFunc("/",func(w http.ResponseWriter, r* http.Request){
        http.ServeFile(w,r,"websocket.html")
    })

    http.ListenAndServe(":8000",nil)
}