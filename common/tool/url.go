package tool

import (
	"net/url"
	"strings"
)

func UrlQueryDecode(q string) (string, string, error) {
	ul, err := url.QueryUnescape(q)
	if err != nil {
		return "", "", err
	}
	ul = strings.Replace(ul, "\\u0026", "&", -1)
	//
	//zhihu
	if strings.Contains(ul, "zhihu") {
		r := strings.Split(ul, "cb=")
		if len(r) == 2 {
			return "zhihu", r[1], nil
		}
	}
	//
	r := strings.Split(ul, "&")
	m := make(map[string]string)
	for idx := 0; idx < len(r); idx++ {
		i := strings.IndexAny(r[idx], "=")
		if i >= 0 {
			m[r[idx][:i]] = r[idx][i+1:]
		}
	}
	return m["source"], ul, nil
}
