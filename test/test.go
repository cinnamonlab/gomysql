package main
import (
	"github.com/cinnamonlab/gomysql"
	"fmt"
	"time"
	"strconv"
)

type testRow struct {
	Id int					`db:"id"`
	Created_at time.Time	`db:"created_at"`
	Updated_at time.Time	`db:"updated_at"`
	Name string				`db:"name"`
	Age float32				`db:"age"`
	Address string			`db:"adress"`
	IsActive int64			`db:"is_active"`
}

func main() {
	config := gomysql.DBConfig{
		Host:"localhost",
		DBName:"test",
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
			rows,err := db.Select("select * from test ")

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

		// test INSERT
		lastInsertId,ierr := db.Insert("insert into test(name,age,adress) values (?,?,?)","ccccc",20,"SDSDSDS")
		if ierr !=nil {
			fmt.Println("Error:"+ ierr.Error())
		} else {
			fmt.Println("Last Insert ID:" + strconv.FormatInt(lastInsertId,10))
		}

		// test Update
		effected,uerr := db.Update("update test set name =? where id=?","xxxxx",lastInsertId)
		if uerr !=nil {
			fmt.Println("Error:"+ uerr.Error())
		} else {
			fmt.Println("Updated row effect:" + strconv.FormatInt(effected,10))
		}

		// test delete
		deffected,derr := db.Delete("DELETE FROM TEST WHERE ID =?",lastInsertId)
		if derr !=nil {
			fmt.Println("Error:"+ derr.Error())
		} else {
			fmt.Println("Deleted row effect:" + strconv.FormatInt(deffected,10))
		}
	}
}
