package Game_Logic

import (
	dataStruct "MeetYou/src/Data_Structure"
	"encoding/json"
	"log"
	"time"
)

func SortGrade(grade []AddGrade){
	for i:=0;i<len(grade);i++{
		for j:=i+1;j<len(grade);j++{
			if grade[i].Grade < grade[j].Grade{
				grade[i],grade[j] = grade[j],grade[i]
			}
		}
	}
}

func judgeStatus(ws *dataStruct.WSConnection){
	for _,v := range ws.WhichTeam.MemberGroup{
		if v.Show == false{
			QuitTeam(v.WSConn,[]byte("{\"userid\":"+ v.UserID+"}"))
		}
		v.Show = false
	}
}


func setHouseOwner(ws *dataStruct.Team){
	var userId string
	for _,v := range ws.MemberGroup{
		userId = v.UserID
	}
	ws.Mutex.Lock()
	ws.MemberGroup[userId].HouseOwner = true
	ws.HouseOwner = userId
	ws.Mutex.Unlock()
}

func sendGrade(ws *dataStruct.Team,path string){
	time.Sleep(1*time.Second)
	ws.SendGrade = false
	rgrade := RAddGrade{
		Path: path,
		Grade:make([]AddGrade,5),
	}
	var i int
	for _,v := range ws.MemberGroup{
		rgrade.Grade[i] = AddGrade{
			UserID:v.UserID,
			UserName:v.UserName,
			Grade:ws.GameGrade[v.UserID],
		}
		i++
	}
	rgrade.Grade = rgrade.Grade[:i]
	SortGrade(rgrade.Grade) //sort the grade
	grade,err := json.Marshal(rgrade)
	if err != nil{
		log.Println(err)
		return
	}
	for _,v := range ws.MemberGroup{
		v.WSConn.OutChan <- grade
	}
}