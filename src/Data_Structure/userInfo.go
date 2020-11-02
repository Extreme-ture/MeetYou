package Data_Structure

import (
	"sync"
	"time"
)

var(
	AllUser  map[string]*User
	AllTeam  map[string]*Team
)

func init(){
	AllUser = make(map[string]*User,1000)
	AllTeam = make(map[string]*Team,1000)
}

type Subject struct{
	Name     string `json:"name"`
	IsAnwser bool   `json:"isanwser"`
	Url      string `json:"url"`
}

type Team struct{
	TeamID            string
	MemberGroup       map[string]*User
	CurrentAnwserID   string
	Mutex             sync.Mutex
	Subjects          [][]*Subject
	HouseOwner        string
	GameStatue        bool
	GameGrade         map[string]int
	SendGrade         bool
}

func(team *Team) TeamDelete(){
	delete(AllTeam,team.TeamID)
}

type User struct{
	WhichTeam      *Team
	ImageUrl       string
	UserName       string
	UserID         string
	WSConn         *WSConnection //websocket连接
	Grade          int
	CurrentGrade   int
	CloseChan      chan bool
	Show           bool
	HouseOwner     bool
	Mutex          sync.Mutex
}

func (user *User) TimeToDelete() {
	select {
	case <-user.CloseChan:
	case <-time.After(3*time.Minute):
		user.WhichTeam = nil
		user.Mutex.Lock()
		delete(AllUser, user.UserID)
		user.Mutex.Unlock()
	}
}