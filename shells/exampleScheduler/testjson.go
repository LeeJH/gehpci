package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type mys struct {
	Name  string `json:"name"`
	Name2 string `json:"name2"`
	name3 string
	Name4 interface{}
}

func main() {
	fmt.Printf("just a test ")
	myss := make([]mys, 0)
	//var myss interface{}
	jsons := `[{"name":"hello"},{"name":"world"}]`
	json.Unmarshal([]byte(jsons), &myss)
	fmt.Printf("rest : %#v \n", myss)
	myss2 := make([]mys, 1)
	myss2[0] = mys{Name: "foo"}
	myss2bit, _ := json.Marshal(myss2)
	fmt.Printf("shoud be : %#v %#v\n", string(myss2bit), jsons)
	rangeStruct(&myss2[0])
}

func rangeStruct(obj interface{}) {
	s := reflect.ValueOf(obj).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !f.CanInterface() {
			continue
		}
		fmt.Printf("%d : %s %s = %v \n", i,
			s.Type().Field(i).Name,
			f.Type(),
			f.Interface())
		jsonb, _ := json.Marshal(f.Interface())
		fmt.Println(s.Type().Field(i).Name + "=" + string(jsonb))

	}
}
