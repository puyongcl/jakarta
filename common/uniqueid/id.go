package uniqueid

import (
	"fmt"
	"jakarta/common/tool"
	"time"
)

// 生成唯一数据库id
func GenDataId() string {
	return tool.Md5ByString(fmt.Sprintf("%d%d", time.Now().UnixNano(), GenId()))
}
