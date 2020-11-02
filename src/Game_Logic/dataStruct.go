package Game_Logic

import dataStruct "MeetYou/src/Data_Structure"

func init(){
	dataStruct.AllHandle = make(map[string]func(*dataStruct.WSConnection,[]byte),15)
	//AllHandle["createTeam"] = CreateTeamInfo
	dataStruct.AllHandle["quitteam"]   = QuitTeam
	dataStruct.AllHandle["gamestart"]  = GameStart
	dataStruct.AllHandle["addsubject"]  = AddSubject
	dataStruct.AllHandle["createteaminfo"]  = CreateTeamInfo
	dataStruct.AllHandle["createuserinfo"]  = CreateUserInfo
	dataStruct.AllHandle["addtoteam"]  = AddToTeam
	dataStruct.AllHandle["gameend"]    = GameEnd
	dataStruct.AllHandle["gradecount"] = GradeCount
	dataStruct.AllHandle["gamecontinue"] = GameContinue
}

type TeamInfo struct{
	InvitationCode  string `json:"invitationcode"`
}

type UserInfo struct {
	Path     string `json:"path"`
	UserName string `json:"username"`
	UserID   string `json:"userid"`
	ImageURL string `json:"imageurl"`
	Rank     int    `json:"rank"`
	Show     bool   `json:"show"`
}

type AllUserInfo struct {
	Path     string `json:"path"`
	Message  bool   `json:"message"`
	AllUser  []UserInfo `json:"alluser"`
}

type QuitInfo struct {
	Path     string `json:"path"`
	UserID    string `json:"userid"`
}

type GameStartInfo struct{
	Path     string `json:"path"`
	UserID    string `json:"userid"`
}

type SubjectInfo struct {
	Path     string `json:"path"`
	UserID     string                      `json:"userid"`
	Subject    [][]*dataStruct.Subject     `json:"subject"`
}

type AddGrade struct {
	Path     string `json:"path"`
	UserID   string  `json:"userid"`
	Grade    int     `json:"grade"`
	UserName string  `json:"username"`
	ImageUrl string  `json:"imageurl"`
	Rank     int  `json:"rank"`
}

type RAddGrade struct {
	Path     string `json:"path"`
	Grade   []AddGrade `json:"grade"`
}
