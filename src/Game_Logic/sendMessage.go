package Game_Logic

import (
	"MeetYou/src/Data_Structure"
	"encoding/json"
	"log"
)

func SendUserInfo(ws *Data_Structure.WSConnection,info *UserInfo){
	ui,err := json.Marshal(info)
	if err != nil{
		log.Println(err)
		return
	}
	ws.OutChan <- ui
}

func SendSubject(ws *Data_Structure.WSConnection,info SubjectInfo){
	ui,err := json.Marshal(info)
	if err != nil{
		log.Println(err)
		return
	}
	ws.OutChan <- ui
}