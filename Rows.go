package gomysql
import (
	"database/sql"
	"reflect"
	"errors"
	"fmt"
)

type Rows struct {
	sql.Rows
	First_one bool
}

func (rows *Rows) ToStruct(dest interface{}) error {

	destMap, err := makeMapable(dest)

	if err != nil {
		return err
	}

	value := reflect.ValueOf(dest)

	direct := reflect.Indirect(value)

	slice, err := baseType(value.Type(), reflect.Slice)

	if err != nil {
		return err
	}

	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := deref(slice.Elem())

	columns, err := rows.Rows.Columns()
	if err != nil {
		return nil
	}

	/*values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}*/

	for rows.Rows.Next() {

		values := make([]interface{}, len(columns))

		vp := reflect.New(base)

		v := reflect.Indirect(vp)

		err=fieldsByTraversal(v, destMap,columns,values,true)

		err = rows.Rows.Scan(values...)
		if err != nil {
			return err
		}

		if isPtr {
			direct.Set(reflect.Append(direct,vp))
		} else {
			direct.Set(reflect.Append(direct,v))
		}

		if(rows.First_one ==true) {
			return nil
		}

	}
	return nil
}

func fieldsByTraversal(v reflect.Value, destMap map[string]reflect.StructField, columns []string, values []interface{}, ptrs bool) error {

	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, name := range columns {
		if fieldMap, ok := destMap[name]; ok {
			f := v.FieldByName(fieldMap.Name)
			if ptrs {
				values[i] = f.Addr().Interface()
			} else {
				values[i] = f.Interface()
			}
		} else {
			values[i] = new(interface{})
		}
	}
	return nil
}

// map Rows objects to Map
func (rows *Rows) ToMap() ([]map[string]interface{}, error) {

	columns, err := rows.Rows.Columns()
	if err != nil {
		return nil,err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	rowMaps := make([]map[string]interface{}, 0)

	for rows.Rows.Next() {
		err = rows.Rows.Scan(values...)
		if err != nil {
			return nil,err
		}

		currRow := make(map[string]interface{})
		for i, name := range columns {
			currRow[name] = *(values[i].(*interface{}))
		}
		// accumulating rowMaps is the easy way out
		rowMaps = append(rowMaps, currRow)
	}

	return rowMaps,nil
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

	for i := 0; i < val.NumField(); i++ {
		//valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		if tag != "" {
			mapobj[tag.Get("json")]=typeField
		}
	}

	return mapobj,nil
}

func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func baseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}
