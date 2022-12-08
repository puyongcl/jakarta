package tool

import (
	"fmt"
	"testing"
)

func Test_abs(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"1",
			args{
				n: -64,
			},
			64,
		},
		{
			"1",
			args{
				n: 64,
			},
			64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.args.n); got != tt.want {
				t.Errorf("abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFmt(t *testing.T) {
	fmt.Printf("%.1f", 12311.1231)
	a := 1
	b := 1
	fmt.Println(a / b)
}
