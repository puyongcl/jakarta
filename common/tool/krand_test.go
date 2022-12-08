package tool

import "testing"

func TestMd5ByString(t *testing.T) {
	s := Md5ByString("每个女孩子都或多或少的有一些容貌焦虑，但是看来你的焦虑有点严重，很影响你的生活，这样的状态持续了有多久了呢？")
	t.Log(s)
}
