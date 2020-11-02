package Game_Logic

import (
	dataStruct "MeetYou/src/Data_Structure"
	conn "MeetYou/src/MySQL_Connection"
	"encoding/json"
	"log"
	"sync"
)

func CreateUserInfo(ws *dataStruct.WSConnection,message []byte){
	ui := new(UserInfo)
	err := json.Unmarshal(message,ui)
	if err != nil{
		log.Println(err)
		return
	}
	//if the user is already exist(when switch the internet or disconnection)
	if _,ok := dataStruct.AllUser[ui.UserID];ok{
		//judge the TimeToDelete(user) is running
		if dataStruct.AllUser[ui.UserID].WSConn == nil{
			dataStruct.AllUser[ui.UserID].CloseChan <- true
		}
		dataStruct.AllUser[ui.UserID].WSConn = ws
		ws.WhichUser = dataStruct.AllUser[ui.UserID]
		if ws.WhichUser.WhichTeam != nil{
			//is add to team
		}
		ui.Rank = ws.WhichUser.Grade
		SendUserInfo(ws,ui)
		return
	}
	userinfo := &dataStruct.User{
		UserName:     ui.UserName,
		ImageUrl:     ui.ImageURL,
		WhichTeam:    nil,
		UserID:       ui.UserID,
		WSConn:       ws,
		Grade:        0,
		CurrentGrade: 0,
		CloseChan:    make(chan bool),
		Mutex:        sync.Mutex{},
	}
	ws.WhichUser = userinfo
	dataStruct.AllUser[ui.UserID] = userinfo
	conn.SelectRow(ui.UserID,&ui.Rank)
	userinfo.Grade = ui.Rank
	SendUserInfo(ws,ui)
}

//cannot prevent create much Team by one user!!!
func CreateTeamInfo(ws *dataStruct.WSConnection,message []byte){
	if ws.WhichUser == nil{
		ws.OutChan <- []byte(`{"path":"createteaminfo","message":"false"}`)
		return
	}
	if ws.WhichTeam != nil{
		ws.OutChan <- []byte(`{"path":"createteaminfo","message":"false"}`)
		return
	}
	ti := new(TeamInfo)
	err := json.Unmarshal(message,ti)
	if err != nil{
		log.Println(err)
		return
	}
	if _,ok := dataStruct.AllTeam[ti.InvitationCode];ok{
		ws.OutChan <- []byte(`{"path":"createteaminfo","message":"false"}`)
		return
	}
	ws.WhichUser.HouseOwner = true  //set the house _owner
	ws.WhichUser.Show = true  //ready to start game
	team := &dataStruct.Team{
		TeamID:          ti.InvitationCode,
		MemberGroup:     make(map[string]*dataStruct.User),
		CurrentAnwserID: "",
		Mutex:           sync.Mutex{},
		Subjects:        make([][]*dataStruct.Subject,0),
		HouseOwner:      ws.WhichUser.UserID,
		GameGrade:       make(map[string]int,5),
	}
	team.MemberGroup[ws.WhichUser.UserID] = ws.WhichUser
	ws.WhichUser.Show = true
	ws.WhichUser.WhichTeam = team
	ws.WhichTeam = team
	dataStruct.AllTeam[team.TeamID] = team
	ws.OutChan <- []byte(`{"path":"createteaminfo","message":"true"}`)
}

func AddToTeam(ws *dataStruct.WSConnection,message []byte){
	if ws.WhichUser == nil{
		return
	}
	if ws.WhichTeam != nil{
		ws.OutChan <- []byte(`{"path":"addtoteam","message":"false"}`)
		return
	}
	ti := new(TeamInfo)
	err := json.Unmarshal(message,ti)
	if err != nil{
		log.Println(err)
		return
	}
	if _,ok := dataStruct.AllTeam[ti.InvitationCode];!ok{
		ws.OutChan <- []byte(`{"path":"addtoteam","message":"false"}`)
		return
	}
	if dataStruct.AllTeam[ti.InvitationCode].GameStatue{
		ws.OutChan <- []byte(`{"path":"addtoteam","message":"false"}`)
		return
	}

	//reject to add to team when member over 5
	if len(dataStruct.AllTeam[ti.InvitationCode].MemberGroup) == 5{
		ws.OutChan <- []byte(`{"path":"addtoteam","message":"false"}`)
		return
	}

	ws.WhichTeam = dataStruct.AllTeam[ti.InvitationCode] //bind ws with team
	ws.WhichUser.WhichTeam = ws.WhichTeam  //bind user with team
	ws.WhichTeam.MemberGroup[ws.WhichUser.UserID] = ws.WhichUser
	ws.WhichUser.Show = true  //ready to start game

	userInfo := make([]UserInfo,5)
	var i int
	for _,v:=range ws.WhichTeam.MemberGroup{
		userInfo[i] = UserInfo{
			Show:     v.Show,
			UserName: v.UserName,
			UserID:   v.UserID,
			ImageURL: v.ImageUrl,
			Rank:     v.Grade,
		}
		i++
	}
	allUser := AllUserInfo{
		Path: "addtoteam",
		Message: true,
		AllUser: userInfo[:i],
	}
	data,err := json.Marshal(allUser)
	if err != nil{
		log.Println(err)
		return
	}
	for _,v := range ws.WhichTeam.MemberGroup{
		v.WSConn.OutChan <- data
	}
}

//send to all user which in the team
func QuitTeam(ws *dataStruct.WSConnection,message []byte){
	qi := new(QuitInfo)
	err := json.Unmarshal(message,qi)
	if err != nil{
		log.Println(err)
		return
	}
	if ws.WhichTeam == nil{
		log.Println("illegal information")
		return
	}
	if ws.WhichUser == nil{
		log.Println("illegal information")
		return
	}
	//just can delete the user in your own team
	if _,ok := ws.WhichTeam.MemberGroup[qi.UserID];!ok{
		log.Println("illegal information")
		return
	}
	if _,ok := dataStruct.AllUser[qi.UserID];!ok{
		log.Println("illegal information")
		return
	}
	if dataStruct.AllUser[qi.UserID].WhichTeam == nil{
		ws.OutChan <- []byte(`{"path":"quitteam","message":"false"}`)
		return
	}
	dataStruct.AllUser[qi.UserID].WSConn.OutChan <- []byte(`{"path":"quitteam","message":"false"}`)
	delete(ws.WhichTeam.MemberGroup,qi.UserID)
	delete(ws.WhichTeam.GameGrade,qi.UserID)
	//if the count of member equal to 0
	if len(ws.WhichTeam.MemberGroup) == 0{
		ws.WhichTeam.TeamDelete()
	}else{
		if ws.WhichUser.HouseOwner{
			setHouseOwner(ws.WhichTeam)  //reselect the homeowner
		}
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
		allUser := AllUserInfo{
			Path: "quitteam",
			Message: true,
			AllUser: userInfo[:i],
		}
		data,err := json.Marshal(allUser)
		if err != nil{
			log.Println(err)
			return
		}
		for _,v := range ws.WhichTeam.MemberGroup{
			v.WSConn.OutChan <- data
		}
	}
	dataStruct.AllUser[qi.UserID].WhichTeam = nil
	dataStruct.AllUser[qi.UserID].WSConn.WhichTeam = nil //break the link between webSocketconn and team
}

func GameStart(ws *dataStruct.WSConnection,message []byte){
	if ws.WhichUser == nil{
		return
	}
	if ws.WhichTeam == nil{
		return
	}
	judgeStatus(ws)
	ws.WhichTeam.GameStatue = true
	for _,value := range ws.WhichTeam.MemberGroup{
		SendSubject(value.WSConn,SubjectInfo{
			Path: "gamestart",
			UserID:  "",
			Subject: ws.WhichTeam.Subjects,
		})
	}
}

func GameEnd(ws *dataStruct.WSConnection,message []byte){
	if ws.WhichTeam == nil{
		return
	}
	if ws.WhichUser == nil{
		return
	}
	//conn.UpdateRow(ws.WhichUser.UserID,ws.WhichUser.Grade)
	sendGrade(ws.WhichTeam,"gameend")
	grade := ws.WhichUser.CurrentGrade
	switch {
	case grade >= 80 :
		conn.UpdateRow(ws.WhichUser.UserID,3)
	case grade >= 65 :
		conn.UpdateRow(ws.WhichUser.UserID,2)
	case grade >= 50 :
		conn.UpdateRow(ws.WhichUser.UserID,1)
	default:

	}
	ws.WhichTeam.GameStatue = false
	ws.WhichUser.CurrentGrade = 0
	if len(ws.WhichTeam.Subjects) > 0{
		ws.WhichTeam.Subjects = make([][]*dataStruct.Subject,0)
	}
}

func AddSubject(ws *dataStruct.WSConnection,message []byte){
	if ws.WhichTeam == nil{
		return
	}
	asi := new(SubjectInfo)
	err := json.Unmarshal(message,asi)
	if err != nil{
		log.Println(err)
		return
	}
	ws.WhichTeam.Mutex.Lock()
	if len(ws.WhichTeam.Subjects) < 20{
		ws.WhichTeam.Subjects = append(ws.WhichTeam.Subjects,asi.Subject...)
	}
	ws.WhichTeam.Mutex.Unlock()

}

func GradeCount(ws *dataStruct.WSConnection,message []byte){
	if ws.WhichUser == nil{
		return
	}
	if ws.WhichTeam == nil{
		return
	}
	ag := new(AddGrade)
	err := json.Unmarshal(message,ag)
	if err != nil{
		log.Println(err)
		return
	}
	ws.WhichUser.CurrentGrade += ag.Grade
	ws.WhichTeam.GameGrade[ws.WhichUser.UserID] = ws.WhichUser.CurrentGrade

	if !ws.WhichTeam.SendGrade {
		ws.WhichTeam.SendGrade = true
		sendGrade(ws.WhichTeam,"gradecount")
	}
}

func GameContinue(ws *dataStruct.WSConnection,message []byte){
	if ws.WhichUser == nil{
		log.Println("not create userInfo")
		return
	}
	if ws.WhichTeam == nil{
		log.Println("not add to team")
		return
	}
	ws.WhichUser.Show = true
	userInfo := make([]UserInfo,5)
	var i int
	for _,v:=range ws.WhichTeam.MemberGroup{
		userInfo[i] = UserInfo{
			Show:     v.Show,
			UserName: v.UserName,
			UserID:   v.UserID,
			ImageURL: v.ImageUrl,
			Rank:     v.Grade,
		}
		i++
	}
	allUser := AllUserInfo{
		Path: "gamecontinue",
		Message: true,
		AllUser: userInfo[:i],
	}
	data,err := json.Marshal(allUser)
	if err != nil{
		log.Println(err)
		return
	}
	for _,v := range ws.WhichTeam.MemberGroup{
		v.WSConn.OutChan <- data
	}
}