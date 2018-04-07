package wxpay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	k_WXPAY_SANDBOX_API_URL    = "https://api.mch.weixin.qq.com/sandboxnew"
	k_WXPAY_PRODUCTION_API_URL = "https://api.mch.weixin.qq.com"
)

type WXPay struct {
	appId        string
	apiKey       string
	mchId        string
	Client       *http.Client
	apiDomain    string
	NotifyURL    string
	isProduction bool
}

func New(appId, apiKey, mchId string, isProduction bool) (client *WXPay) {
	client = &WXPay{}
	client.appId = appId
	client.mchId = mchId
	client.apiKey = apiKey
	client.Client = http.DefaultClient
	client.isProduction = isProduction
	if isProduction {
		client.apiDomain = k_WXPAY_PRODUCTION_API_URL
	} else {
		client.apiDomain = k_WXPAY_SANDBOX_API_URL
	}
	return client
}

func (this *WXPay) URLValues(param WXPayParam, key string) (value url.Values, err error) {
	var p = param.Params()
	p.Set("appid", this.appId)
	p.Set("mch_id", this.mchId)
	p.Set("nonce_str", getNonceStr())

	if _, ok := p["notify_url"]; ok == false {
		if len(this.NotifyURL) > 0 {
			p.Set("notify_url", this.NotifyURL)
		}
	}

	var sign = signMD5(p, key)
	p.Set("sign", sign)
	return p, nil
}

func (this *WXPay) doRequest(method, url string, param WXPayParam, results interface{}) (err error) {
	key, err := this.getKey()
	if err != nil {
		return err
	}

	p, err := this.URLValues(param, key)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, strings.NewReader(urlValueToXML(p)))
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

	if _, err := verifyResponseData(data, key); err != nil {
		return err
	}

	err = xml.Unmarshal(data, results)

	return err
}

func (this *WXPay) DoRequest(method, url string, param WXPayParam, results interface{}) (err error) {
	return this.doRequest(method, url, param, results)
}

func (this *WXPay) getKey() (key string, err error) {
	if this.isProduction == false {
		key, err = this.getSignKey(this.apiKey)
		if err != nil {
			return "", err
		}
	} else {
		key = this.apiKey
	}
	return key, err
}

func (this *WXPay) getSignKey(apiKey string) (key string, err error) {
	var p = make(url.Values)
	p.Set("mch_id", this.mchId)
	p.Set("nonce_str", getNonceStr())

	var sign = signMD5(p, apiKey)
	p.Set("sign", sign)

	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/sandboxnew/pay/getsignkey", strings.NewReader(urlValueToXML(p)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	resp, err := this.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var signKey *GetSignKeyResp
	if err = xml.Unmarshal(data, &signKey); err != nil {
		return "", err
	}
	return signKey.SandboxSignKey, nil
}

func (this *WXPay) BuildAPI(paths ...string) string {
	var path = this.apiDomain
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			if strings.HasSuffix(path, "/") {
				path = path + p
			} else {
				if strings.HasPrefix(p, "/") {
					path = path + p
				} else {
					path = path + "/" + p
				}
			}
		}
	}
	return path
}

func urlValueToXML(m url.Values) string {
	var xmlBuffer = &bytes.Buffer{}
	xmlBuffer.WriteString("<xml>")

	for key := range m {
		var value = m.Get(key)
		if key == "total_fee" || key == "refund_fee" || key == "execute_time_" {
			xmlBuffer.WriteString("<" + key + ">" + value + "</" + key + ">")
		} else {
			xmlBuffer.WriteString("<" + key + "><![CDATA[" + value + "]]></" + key + ">")
		}
	}
	xmlBuffer.WriteString("</xml>")
	return xmlBuffer.String()
}

func signMD5(param url.Values, key string) (sign string) {
	var keys = make([]string, 0, 0)
	for key := range param {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = param.Get(key)
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	if key != "" {
		pList = append(pList, "key="+key)
	}

	var src = strings.Join(pList, "&")
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	cipherStr := md5Ctx.Sum(nil)

	sign = strings.ToUpper(hex.EncodeToString(cipherStr))
	return sign
}

// TODO 用于验证 notify url
func verifySign(data url.Values, key string) (ok bool, err error) {
	return ok, err
}

func verifyResponseData(data []byte, key string) (ok bool, err error) {
	var param = make(XMLMap)
	err = xml.Unmarshal(data, &param)
	if err != nil {
		return false, err
	}

	// 处理错误信息
	var code = param.Get("return_code")
	if code == K_RETURN_CODE_FAIL {
		var msg = param.Get("return_msg")
		return false, errors.New(msg)
	}

	code = param.Get("result_code")
	if code == K_RETURN_CODE_FAIL {
		var msg = param.Get("err_code_des")
		return false, errors.New(msg)
	}

	// 验证签名
	var sign = param.Get("sign")
	delete(param, "sign")
	if sign == "" {
		return false, errors.New("签名验证失败")
	}

	var sign2 = signMD5(url.Values(param), key)
	if sign == sign2 {
		return true, nil
	}
	return false, errors.New("签名验证失败")
}

func getNonceStr() (nonceStr string) {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 32; i++ {
		idx := rand.Intn(len(chars) - 1)
		nonceStr += chars[idx : idx+1]
	}
	return nonceStr
}
