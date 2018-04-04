package wxpay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
)

type WXPay struct {
	appId     string
	apiKey    string
	mchId     string
	Client    *http.Client
	NotifyURL string
}

func New(appId, apiKey, mchId string) (client *WXPay) {
	client = &WXPay{}
	client.appId = appId
	client.mchId = mchId
	client.apiKey = apiKey
	client.Client = http.DefaultClient
	return client
}

func (this *WXPay) doRequest(method, url string, param WXPayParam, results interface{}) (err error) {
	var p = param.Params()
	p["appid"] = this.appId
	p["mch_id"] = this.mchId
	p["nonce_str"] = getNonceStr()
	if _, ok := p["notify_url"]; ok == false {
		if len(this.NotifyURL) > 0 {
			p["notify_url"] = this.NotifyURL
		}
	}

	var sign = signMD5(p, this.apiKey)
	p["sign"] = sign

	req, err := http.NewRequest(method, url, strings.NewReader(mapToXML(p)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	resp, err := this.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, results)

	return err
}

func (this *WXPay) DoRequest(method, url string, param WXPayParam, results interface{}) (err error) {
	return this.doRequest(method, url, param, results)
}

func mapToXML(m map[string]interface{}) string {
	var xmlBuffer = &bytes.Buffer{}
	xmlBuffer.WriteString("<xml>")

	for key, value := range m {
		var value = fmt.Sprintf("%v", value)
		if key == "total_fee" || key == "refund_fee" || key == "execute_time_" {
			xmlBuffer.WriteString("<" + key + ">" + value + "</" + key + ">")
		} else {
			xmlBuffer.WriteString("<" + key + "><![CDATA[" + value + "]]></" + key + ">")
		}
	}
	xmlBuffer.WriteString("</xml>")
	return xmlBuffer.String()
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
		pList = append(pList, "key="+apiKey)
	}

	var src = strings.Join(pList, "&")
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	cipherStr := md5Ctx.Sum(nil)

	sign = strings.ToUpper(hex.EncodeToString(cipherStr))
	return sign
}

func getNonceStr() (nonceStr string) {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 32; i++ {
		idx := rand.Intn(len(chars) - 1)
		nonceStr += chars[idx : idx+1]
	}
	return nonceStr
}
