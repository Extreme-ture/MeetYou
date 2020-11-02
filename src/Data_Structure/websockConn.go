package Data_Structure

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var(
	AllHandle map[string]func(*WSConnection,[]byte)
)


// http升级websocket协议的配置
var WSUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


// 客户端连接
type WSConnection struct {
	WhichTeam   *Team
	WhichUser   *User
	WsSocket    *websocket.Conn // 底层websocket
	InChan      chan []byte	// 读队列
	OutChan     chan []byte // 写队列
	mutex       sync.Mutex	// 避免重复关闭管道
	IsClosed    bool
	CloseChan   chan byte  // 关闭通知
}

func(ws *WSConnection) ReadF(){
	for{
		_,data,err := ws.WsSocket.ReadMessage()
		if err != nil{
			goto erro
		}

		// 放入请求队列
		select {
		case ws.InChan <- data:
		case <- ws.CloseChan:
			goto closed
		}
	}

erro:
	ws.WSClose()
closed:
}

func(ws *WSConnection) WriteT(){
	for {
		select {
		// 取一个应答
		case msg := <- ws.OutChan:
			// 写给websocket
			if err := ws.WsSocket.WriteMessage(1, msg); err != nil {
				goto error
			}
		case <- ws.CloseChan:
			goto closed
		}
	}
error:
	ws.WSClose()
closed:
}

//handle the path
type handPath struct {
	Path  string `json:"path"`
}

func(ws *WSConnection) HandleMessage(){
	for{
		select {
			case message := <-ws.InChan:
				path := new(handPath)
				err := json.Unmarshal(message,path)
				if err != nil{
					log.Println(ws.WhichUser.UserID," ",err)
					continue
				}
				//select the func to solve
				if v,ok := AllHandle[path.Path];ok{
					v(ws,message)
				}
		    case <-ws.CloseChan:
			goto ERR
		}
	}
ERR:
}

func (ws *WSConnection)WSClose() {
	_ = ws.WsSocket.Close()

	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	if !ws.IsClosed {
		ws.IsClosed = true
		close(ws.CloseChan)
		if ws.WhichUser == nil{
			return
		}else if ws == ws.WhichUser.WSConn{
			QuitTeam(ws)
		}else{
			ws.WhichUser = nil
			ws.WhichTeam = nil
		}
	}
}

type UserInfo struct {
	UserName string `json:"username"`
	UserID   string `json:"userid"`
	ImageURL string `json:"imageurl"`
	Rank     int    `json:"rank"`
}

func QuitTeam(ws *WSConnection){
	if ws.WhichTeam != nil{
		delete(ws.WhichTeam.MemberGroup,ws.WhichUser.UserID) //quit team
		userInfo := make([]UserInfo,5)
		var i int
		for _,v:=range ws.WhichTeam.MemberGroup{
			userInfo[i] = UserInfo{
				UserName: v.UserName,
				UserID:   v.UserID,
				ImageURL: v.ImageUrl,
				Rank:     v.Grade,
			}
			i++
		}
		data,err := json.Marshal(userInfo)
		if err != nil{
			log.Println(err)
			return
		}
		for _,v := range ws.WhichTeam.MemberGroup{
			v.WSConn.OutChan <- data
		}

		if len(ws.WhichTeam.MemberGroup) == 0{
			ws.WhichTeam.TeamDelete()
		}
	}
	ws.WhichTeam = nil //break the link between webSocketconn and team
	ws.WhichUser.WSConn = nil
	user := ws.WhichUser
	ws.WhichUser = nil
	go user.TimeToDelete()
}
