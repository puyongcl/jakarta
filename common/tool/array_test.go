package tool

import (
	"fmt"
	"github.com/jinzhu/copier"
	"testing"
)

func TestCombineStringArray(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"c", "d", "e"}
	fmt.Println(CombineStringArray(a, b))
}

func TestGetIntArrayString(t *testing.T) {
	a := []int64{1, 2, 3, 4}
	fmt.Println(GetIntArrayString(a))
}

func TestCopy(t *testing.T) {
	type arr1 struct {
		Arr []int64
	}
	type arr2 struct {
		Arr []int64
	}
	a := arr1{Arr: make([]int64, 0)}
	var idx int64
	for ; idx < 4; idx++ {
		b := arr2{Arr: []int64{idx + 1, idx + 2, idx + 3, idx + 4}}
		_ = copier.Copy(&a, &b)
		fmt.Println(a)
	}
}

func TestAppend(t *testing.T) {
	var l []string
	l = []string{"1", "2", "3"}
	l = append(l, "")
	var c int64
	for k, _ := range l {
		fmt.Println(l[k])
		c++
	}
	fmt.Println(l, c)
}

func TestSlice(t *testing.T) {
	var a []int64
	a = []int64{1, 2, 3, 4, 5}
	a = a[1:]
	fmt.Println(a)
}
