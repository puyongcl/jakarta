package tool

import (
	"fmt"
	"testing"
	"time"
)

func TestGetAge2(t *testing.T) {
	d, _ := time.Parse(DateLayout, "1985-02-09")
	fmt.Println(GetAge2(d))
}
