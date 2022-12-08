package listener

import (
	"fmt"
	"testing"
)

func TestGetBanner(t *testing.T) {
	got := GetBanner()
	fmt.Println(got)
}
