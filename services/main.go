package main

import (
	"errors"
	"fmt"
	"reflect"
)

type User struct {
	Username string
	Email string
	Balance float64
}

func Copy(src interface{}, tgt interface{}) error {
	var tv = reflect.ValueOf(tgt) // reflect.value
	if tv.Kind() != reflect.Ptr || tv.Elem().Kind() != reflect.Struct {
		return errors.New("Value not pointer of a struct")
	}

	// starts the copying process
	var te = tv.Elem()
	var sv = reflect.ValueOf(src) 
	
	for i:=0; i < sv.NumField(); i++ {
		var sf = sv.Field(i)
		var tf = te.Field(i)

		if reflect.Zero(tf.Type()).Interface() == tf.Interface() {
			tf.Set(sf)
		}
	}

	return nil
}

func main(){
	var u = &User{Username: "some"}
	var t = User{"doe","marios@gmail.com", 0.10}

	fmt.Println(Copy(t,u))
	fmt.Println(u)

}