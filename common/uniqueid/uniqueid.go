package uniqueid

import (
	"github.com/sony/sonyflake"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime:      time.Now(),
		MachineID:      nil,
		CheckMachineID: nil,
	})
	if flake == nil {
		panic("flake id not created")
	}
}

func GenId() int64 {
	id, err := flake.NextID()
	if err != nil {
		logx.Severef("flake NextID failed with %s \n", err)
		return time.Now().UnixNano()
	}

	return int64(id)
}
