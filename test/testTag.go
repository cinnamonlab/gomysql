package main

import (
	"fmt"
	"reflect"
	"github.com/cinnamonlab/gomysql"
	"errors"
)

type Foo struct {
	FirstName string `tag_name:"tag 1"`
	LastName  string `tag_name:"tag 2"`
	Age       int		`tag_name:"tag 3"`
}

func (f *Foo) reflect( dest interface{}) error {

	value := reflect.ValueOf(dest)

	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	slice, err := gomysql.BaseType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	base := gomysql.Deref(slice.Elem())

	destStruct := reflect.New(base)

	val := reflect.Indirect(destStruct)

	mapobj := make(map[string]reflect.StructField)

	for i := 0; i < val.NumField(); i++ {
		//valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		if tag != "" {
			mapobj[tag.Get("tag_name")]=typeField
		}
	}

	fmt.Println(mapobj)

	return nil
}

func main() {
	f := &Foo{
		FirstName: "Drew",
		LastName:  "Olson",
		Age:       30,
	}
	pps := []Foo{}

	f.reflect(&pps)
}