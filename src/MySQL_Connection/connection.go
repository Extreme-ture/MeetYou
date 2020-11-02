package MySQL_Connection

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const (
	userName = "root"
	password = "Hz1966.."
	ip = "127.0.0.1"
	port = "3306"
	DBName = "MeetYou"
)

var(
	DB *sql.DB
)

func Connect(){
	//path := strings.Join([]string{userName, ":", password, "@tcp(",ip, ":", port, ")/", DBName, "?charset=utf8mb4:"}, "")
	db,err := sql.Open("mysql","root:Hz1966..@tcp(localhost:3306)/MeetYou?charset=utf8mb4")
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil{
		log.Println(err)
		os.Exit(1)
	}
	DB = db
	log.Println("Open MySql Success")
}

func SelectRow(userID string,Rank *int){
	stmt,_ := DB.Prepare("SELECT level FROM user WHERE userID = ?")
	defer stmt.Close()
	err := stmt.QueryRow(userID).Scan(Rank)
	if err == sql.ErrNoRows{
		InsertRow(userID,*Rank)
	}else{
		if err != nil{
			log.Println(err)
		}
	}
}

func InsertRow(userID string,Rank int){
	stmt, _ := DB.Prepare("INSERT INTO user (userID,level) VALUES (?,?)")
	defer stmt.Close()

	_, err := stmt.Exec(userID,Rank)
	if err != nil {
		log.Println(err)
	}
}

func UpdateRow(userID string,Rank int){
	stmt,_:= DB.Prepare("UPDATE user SET level=? where userID=?")
	defer stmt.Close()
	_,err := stmt.Exec(Rank,userID)
	if err != nil{
		log.Println(err)
	}
}

func DeleteRow(userID string){
	stmt,_ := DB.Prepare("DELETE FROM user WHERE userID=?")
	defer stmt.Close()

	_,err := stmt.Exec(userID)
	if err != nil{
		log.Println(err)
	}
}
