package Internet_interface

import (
	dataStruct "MeetYou/src/Data_Structure"
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func InitLog(){
	logFile, logErr := os.OpenFile("/home/MeetYou/src/myLog/mylog.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Start(){
	http.HandleFunc("/upgrade",upgrade)
	http.HandleFunc("/seekTeam",seekTeam)
	err := http.ListenAndServe("0.0.0.0:8080",nil)
	if err != nil{
		log.Println(err)
	}
}

func seekTeam(w http.ResponseWriter,r *http.Request){
	data,err := json.Marshal(dataStruct.AllTeam)
	if err != nil{
		log.Println(err)
	}
	w.Write(data)
}

func upgrade(w http.ResponseWriter,r *http.Request){
	wsSocket, err := dataStruct.WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	wsConn := &dataStruct.WSConnection{
		WhichTeam: nil,
		WhichUser: nil,
		WsSocket:  wsSocket,
		InChan:    make(chan []byte,20),
		OutChan:   make(chan []byte,20),
		CloseChan: make(chan byte,2),
	}

	//run the read func
	go wsConn.ReadF()
	//run the write func
	go wsConn.WriteT()
	//run the handle func
	go wsConn.HandleMessage()
}
