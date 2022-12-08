package userkey

// 图片文件 https://jakarta-13178454693.cos.ap-guangzhou.myqcloud.com/default/1.jpg
var userAvatar = []string{
	"default/1.jpg",
	"default/2.jpg",
	"default/3.jpg",
	"default/4.jpg",
	"default/5.jpg",
	"default/6.jpg",
	"default/7.jpg",
	"default/8.jpg",
	"default/9.jpg",
}

var avatarListLen int64 = int64(len(userAvatar))

func GetDefaultAvatar(uid int64) string {
	idx := uid % avatarListLen

	return userAvatar[idx]
}
