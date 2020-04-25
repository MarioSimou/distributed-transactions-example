package main

import (
	"fmt"
	"reflect"
)

func copy(src interface{}, tgt interface{}) error {
	var tv = reflect.ValueOf(tgt)
	if tv.Kind() != reflect.Ptr || tv.Elem().Kind() != reflect.Struct {
		return  fmt.Errorf("Error: 'Target is not a pointer of struct'\n")
	}

	var sv = reflect.ValueOf(src)
	var te = tv.Elem()
	for i := 0; i < sv.NumField(); i++ {
		var sf = sv.Field(i)
		var tf = te.Field(i)

		if reflect.Zero(tf.Type()).Interface() == tf.Interface() {
			tf.Set(sf)
		}
	}

	return nil
}

type product struct {
	Id int
	Name string
	Price float64
	Quantity int
}

func main(){
	var src = product{1,"some",0.10, 10}
	var tgt = product{Id: 2}

	fmt.Println(copy(src,&tgt))
	fmt.Println(tgt)
}