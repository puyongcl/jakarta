package uniqueid

import (
	"fmt"
	"testing"
	"time"
)

func TestGenSn(t *testing.T) {
	fmt.Println(GenSn(SnPrefixTextChatOrderId))
	fmt.Println(GenDataId())
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().UnixNano())
}
