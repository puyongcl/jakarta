package tim

import (
	"jakarta/common/key/tencentcloudkey"
	"jakarta/common/xerr"
	"strconv"
)

const updateProfileUrl = "/v4/profile/portrait_set"

type Pair struct {
	Tag   string `json:"Tag"`
	Value string `json:"Value"`
}

type UpdateProfile struct {
	FromAccount string `json:"From_Account"`
	ProfileItem []Pair `json:"ProfileItem"`
}

// 更新IM资料
func (ra *RestApiClient) UpdateProfile(uid int64, nick, avatar string) error {
	req := UpdateProfile{
		FromAccount: strconv.FormatInt(uid, 10),
		ProfileItem: make([]Pair, 0),
	}

	if nick != "" {
		req.ProfileItem = append(req.ProfileItem, Pair{
			Tag:   "Tag_Profile_IM_Nick",
			Value: nick,
		})
	}
	if avatar != "" {
		req.ProfileItem = append(req.ProfileItem, Pair{
			Tag:   "Tag_Profile_IM_Image",
			Value: tencentcloudkey.CDNBasePath + avatar,
		})
	}

	if len(req.ProfileItem) <= 0 {
		return nil
	}

	var rsp CommonResp
	resp, err := ra.Client.R().SetBody(&req).SetResult(&rsp).Post(ra.GetTimRestApiUrl(updateProfileUrl))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, resp.Status())
	}
	if rsp.ActionStatus == "OK" {
		return nil
	}
	return xerr.NewGrpcErrCodeMsg(rsp.ErrorCode, rsp.ErrorInfo)
}
