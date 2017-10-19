package wxpay

import (
	"fmt"
	"sort"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

type Param interface {
	Params() map[string]interface{}
}

func mapToXML(m map[string]interface{}) (xml string) {
	xml = "<xml>"
	for k, v := range m {
		var value = fmt.Sprintf("%v", v)
		if k == "total_fee" || k == "refund_fee" || k == "execute_time_" {
			xml += "<" + k + ">" + value + "</" + k + ">"
		} else {
			xml += "<" + k + "><![CDATA[" + value + "]]></" + k + ">"
		}
	}
	xml += "</xml>"
	return xml
}

func signMD5(param map[string]interface{}, apiKey string) (sign string) {
	var keys = make([]string, 0, 0)
	for key := range param {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = fmt.Sprintf("%v", param[key])
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	if apiKey != "" {
		pList = append(pList, "key="+ apiKey)
	}

	var src = strings.Join(pList, "&")
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	cipherStr := md5Ctx.Sum(nil)

	sign = strings.ToUpper(hex.EncodeToString(cipherStr))
	return sign
}

func getNonceStr() (nonceStr string) {
	chars := "ABCDEFGIHJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < 32; i++ {
		idx := rand.Intn(len(chars) - 1)
		nonceStr += chars[idx : idx+1]
	}
	return nonceStr
}