package tool

import (
	"strconv"
	"strings"
)

func IsStringArrayExist(in string, array []string) bool {
	if len(array) <= 0 {
		return false
	}
	for k, _ := range array {
		if array[k] == in {
			return true
		}
	}
	return false
}

func IsInt64ArrayExist(in int64, array []int64) bool {
	if len(array) <= 0 {
		return false
	}
	for k, _ := range array {
		if array[k] == in {
			return true
		}
	}
	return false
}

func CombineStringArray(a []string, b []string) []string {
	check := make(map[string]struct{})
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = struct{}{}
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}

func IsEqualArrayInt64(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	ma := make(map[int64]struct{})
	for _, v := range a {
		ma[v] = struct{}{}
	}
	for _, v := range b {
		_, ok := ma[v]
		if !ok {
			return false
		}
	}
	return true
}

func IsEqualArrayInt64Order(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := 0; idx < len(a); idx++ {
		if a[idx] != b[idx] {
			return false
		}
	}
	return true
}

func IsEqualArrayString(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	ma := make(map[string]struct{})
	for _, v := range a {
		ma[v] = struct{}{}
	}
	for _, v := range b {
		_, ok := ma[v]
		if !ok {
			return false
		}
	}
	return true
}

func GetIntArrayString(a []int64) string {
	if len(a) <= 0 {
		return ""
	}
	b := make([]string, len(a))
	for idx := 0; idx < len(a); idx++ {
		b[idx] = strconv.FormatInt(a[idx], 10)
	}
	return strings.Join(b, "#")
}
