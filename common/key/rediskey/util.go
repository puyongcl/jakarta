package rediskey

func GetInterfaceArray(in []string) []interface{} {
	temp := make([]interface{}, 0)
	for k, _ := range in {
		temp = append(temp, in[k])
	}
	return temp
}
