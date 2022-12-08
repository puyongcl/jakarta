package tool

import (
	"reflect"
)

// 获取结构体字段名
func GetStructFieldName(d interface{}) []string {
	t := reflect.TypeOf(d)
	//v := reflect.ValueOf(d)
	names := make([]string, 0) // name
	//tags := make([]string, 0)
	//types := make([]string, 0)
	//valuse := make([]interface{}, 0)
	for k := 0; k < t.NumField(); k++ {
		names = append(names, t.Field(k).Name)
		//valuse = append(valuse, v.Field(k).Interface())
		//types = append(types, t.Field(k).Type.Name())
		//tags = append(tags, t.Field(k).Tag.Get("json"))
	}
	//fmt.Println(names)
	//fmt.Println(tags)
	//fmt.Println(types)
	//fmt.Println(valuse)
	return names
}
