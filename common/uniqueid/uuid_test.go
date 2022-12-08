package uniqueid

import (
	"fmt"
	"testing"
)

func TestGenUuid(t *testing.T) {
	got := GenUuid()
	fmt.Println(got)
}
