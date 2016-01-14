package main
import (
	"github.com/cinnamonlab/gomysql"
	"fmt"
	"time"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

type testRow struct {
	Id 			string 			`json:"id"`
	Name 		string 			`json:"name"`
	FirstName 	string 			`json:"first_name"`
	LastName 	string 			`json:"last_name"`
	Avatar 		string 			`json:"avatar"`
	Email 		sql.NullString 	`json:"email"`
	Phone		sql.NullString  `json:"phone"`
	Language 	string 			`json:"language"`
	AccessToken string 			`json:"access_token"`

	Password 	string 			`json:"pass_word"`
	UserName  	string 			`json:"user_name"`

	CreatedAt 	time.Time	 	`json:"created_at"`
	UpdatedAt 	time.Time 		`json:"updated_at"`
	EmailNotiFlag bool 			`json:"email_noti_flg"`

	ScannedAt 	mysql.NullTime  `json:"scanned_at"`
	LastAction 	mysql.NullTime 	`json:"last_action"`
	LastLogin  	mysql.NullTime 	`json:"last_login"`
}

func main() {
	config := gomysql.DBConfig{
		Host:"localhost",
		DBName:"tuya",
		User:"root",
		Password:"",
		Charset:"utf8",
		Port:"3306",
	}
	db,err := gomysql.NewConnection(&config)

	if err !=nil {
		fmt.Println("Can not open Mysql connection")
	} else {
		// test select
			rows,err := db.Select("select * from user ")

		if err !=nil {
			fmt.Println("Error:"+ err.Error())
		} else {
			defer rows.Rows.Close()

			test :=[]testRow{}

			merr := rows.ToStruct(&test)

			if merr !=nil {
				fmt.Println(merr.Error());
			} else
			{
				fmt.Println(test)
			}
		}
	}
}
