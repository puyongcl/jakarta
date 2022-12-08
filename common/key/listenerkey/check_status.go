package listenerkey

import "strings"

// XXX审核状态
const (
	CheckStatusFirstApplyEdit     = 1 // 首次申请 编辑中
	CheckStatusFirstApplyChecking = 2 // 首次申请 提交审核 待审核
	CheckStatusFirstApplyRefuse   = 3 // 首次申请 审核拒绝
	CheckStatusFirstApplyPass     = 4 // 首次申请 审核通过
	CheckStatusEditWaitChecking   = 5 // 修改资料 待审核
	CheckStatusEditRefuse         = 6 // 审核拒绝 管理后台有审核项不通过则为拒绝，通过的项更新
	CheckStatusEditPass           = 7 // 审核通过 全部审核通过
)

// 首次申请审核状态
var FirstApplyListenerCheckStatus = []int64{CheckStatusFirstApplyEdit, CheckStatusFirstApplyChecking, CheckStatusFirstApplyRefuse}

var FirstApplyEditListenerCheckStatus = []int64{CheckStatusFirstApplyEdit, CheckStatusFirstApplyChecking}

// XXX审核状态
var ListenerCheckStatus = []int64{CheckStatusFirstApplyPass, CheckStatusEditWaitChecking, CheckStatusEditRefuse, CheckStatusEditPass}

func GetCheckFieldText(a []string) string {
	if len(a) <= 0 {
		return ""
	}
	b := make([]string, 0)
	for idx := 0; idx < len(a); idx++ {
		r, ok := ListenerProfileCheckFieldText[a[idx]]
		if !ok {
			continue
		}
		b = append(b, r)
	}

	return strings.Join(b, ",")
}
