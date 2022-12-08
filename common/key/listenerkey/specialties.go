package listenerkey

// 擅长领域 例子：一级领域：1-99 二级领域：101-9999 前两位数代表对应一级领域 后两位数表示二级领域
var Specialties = map[int64]string{
	1:    "",
	101:  "",
	102:  "",
	103:  "",
	104:  "",
	105:  "",
	106:  "",
	107:  "",
	108:  "",
	109:  "",
	110:  "",
	111:  "",
	2:    "",
	201:  "",
	202:  "",
	203:  "",
	204:  "",
	205:  "",
	206:  "",
	207:  "",
	208:  "",
	209:  "",
	210:  "",
	211:  "",
	212:  "",
	213:  "",
	214:  "",
	215:  "",
	3:    "",
	301:  "",
	302:  "",
	303:  "",
	304:  "",
	305:  "",
	306:  "",
	307:  "",
	308:  "",
	309:  "",
	310:  "",
	311:  "",
	312:  "",
	313:  "",
	4:    "",
	401:  "",
	402:  "",
	403:  "",
	404:  "",
	405:  "",
	406:  "",
	407:  "",
	408:  "",
	409:  "",
	410:  "",
	411:  "",
	412:  "",
	5:    "",
	501:  "",
	6:    "",
	601:  "",
	602:  "",
	603:  "",
	604:  "",
	605:  "",
	606:  "",
	607:  "",
	608:  "",
	609:  "",
	610:  "",
	611:  "",
	612:  "",
	613:  "",
	614:  "",
	7:    "",
	701:  "",
	702:  "",
	703:  "",
	704:  "",
	705:  "",
	8:    "",
	801:  "",
	802:  "",
	803:  "",
	804:  "",
	805:  "",
	806:  "",
	807:  "",
	808:  "",
	809:  "",
	810:  "",
	9:    "",
	901:  "",
	902:  "",
	903:  "",
	904:  "",
	905:  "",
	906:  "",
	907:  "",
	10:   "",
	1001: "",
	1002: "",
	1003: "",
	1004: "",
	1005: "",
	1006: "",
	1007: "",
	1008: "",
	1009: "",
	1010: "",
}

var SpecialtiesPic = []string{"icon/home-icon-1-2.png", "icon/home-icon-2-2.png", "icon/home-icon-3-2.png", "icon/home-icon-4-2.png", "icon/home-icon-5-2.png", "icon/home-icon-6-2.png", "icon/home-icon-7-2.png", "icon/home-icon-8-2.png", "icon/home-icon-9-2.png", "icon/home-icon-10-2.png"}

var SpecialtiesLevelOneId = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

var SpecialtiesLevelTwoId = [][]int64{
	{101,
		102,
		103,
		104,
		105,
		106,
		107,
		108,
		109,
		110,
		111},
	{201,
		202,
		203,
		204,
		205,
		206,
		207,
		208,
		209,
		210,
		211,
		212,
		213,
		214,
		215}, {
		301,
		302,
		303,
		304,
		305,
		306,
		307,
		308,
		309,
		310,
		311,
		312,
		313},
	{
		401,
		402,
		403,
		404,
		405,
		406,
		407,
		408,
		409,
		410,
		411,
		412},
	{
		501},
	{
		601,
		602,
		603,
		604,
		605,
		606,
		607,
		608,
		609,
		610,
		611,
		612,
		613,
		614},
	{
		701,
		702,
		703,
		704,
		705},
	{
		801,
		802,
		803,
		804,
		805,
		806,
		807,
		808,
		809,
		810},
	{
		901,
		902,
		903,
		904,
		905,
		906,
		907},
	{
		1001,
		1002,
		1003,
		1004,
		1005,
		1006,
		1007,
		1008,
		1009,
		1010},
}