package internal

import (
	"errors"
	"reflect"
)

func copy(src interface{}, tgt interface{}) error {
	var tv = reflect.ValueOf(tgt) // reflect.value
	if tv.Kind() != reflect.Ptr || tv.Elem().Kind() != reflect.Struct {
		return errors.New("Value not pointer of a struct")
	}

	var te = tv.Elem()
	var sv = reflect.ValueOf(src) 
	
	for i:=0; i < sv.NumField(); i++ {
		var sf = sv.Field(i)
		var tf = te.Field(i)

		// check the value of the respect data type empty value
		if reflect.Zero(tf.Type()).Interface() == tf.Interface() {
			tf.Set(sf)
		}
	}
	return nil
}