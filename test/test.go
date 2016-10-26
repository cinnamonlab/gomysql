package main
import (
	"github.com/cinnamonlab/gomysql"
	"fmt"
	"reflect"
	"errors"
)

type FirstTable struct {
	FirstTableId int `json:"idfirst_table"`
	FirstTableCol string `json:"first_tablecol"`
	FirstTableCol1 string `json:"first_tablecol1"`
}

type SecondTable struct {
	SecondTableId int `json:"idsecond_table"`
	SecondTableCol string `json:"second_tablecol"`
	SecondTableCol1 string `json:"second_tablecol1"`
}

type MixTable struct {
	SecondTable
	FirstTable
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
			rows,err := db.Select("select  * from first_table inner join second_table on second_table.first_table_id = first_table.idfirst_table")

		if err !=nil {
			fmt.Println("Error:"+ err.Error())
		} else {
			defer rows.Rows.Close()

			test :=[]MixTable{}

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

func makeMapable(rowsSlicePtr interface{}) (map[string]reflect.StructField, error) {

	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))

	if sliceValue.Kind() != reflect.Slice {
		return nil,errors.New("needs a pointer to a slice")
	}

	sliceElementType := sliceValue.Type().Elem()
	st := reflect.New(sliceElementType)

	val := reflect.Indirect(st)

	mapobj := make(map[string]reflect.StructField)

	getMappableStructure(val, &mapobj)

	return mapobj,nil
}

func getMappableStructure(obj reflect.Value, resultMap *map[string]reflect.StructField) {
	for i := 0; i < obj.NumField(); i++ {
		typeField := obj.Type().Field(i)
		if typeField.Type.Kind() == reflect.Struct {
			newObject := reflect.New(typeField.Type)

			newVal := reflect.Indirect(newObject)

			getMappableStructure(newVal,resultMap)
		}
		tag := typeField.Tag

		if tag != "" {
			(*resultMap)[tag.Get("json")] = typeField
		}
	}
}
