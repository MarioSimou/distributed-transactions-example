package internal

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
)

func HandleCORS(c *gin.Context){
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type: application/json")
	c.Next()
}

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