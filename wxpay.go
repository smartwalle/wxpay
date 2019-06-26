package wxpay

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

type Client struct {
	appId        string
	apiKey       string
	mchId        string
	Client       *http.Client
	tlsClient    *http.Client
	apiDomain    string
	NotifyURL    string
	isProduction bool
}

func New(appId, apiKey, mchId string, isProduction bool) (client *Client) {
	client = &Client{}
	client.appId = appId
	client.mchId = mchId
	client.apiKey = apiKey
	client.Client = http.DefaultClient
	client.isProduction = isProduction
	if isProduction {
		client.apiDomain = kProductionURL
	} else {
		client.apiDomain = kSandboxURL
	}
	return client
}

func initTLSClient(cert []byte, password string) (tlsClient *http.Client, err error) {
	if len(cert) > 0 {
		cert, err := pkcs12ToPem(cert, password)
		if err != nil {
			return nil, err
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		transport := &http.Transport{
			TLSClientConfig:    config,
			DisableCompression: true,
		}

		tlsClient = &http.Client{Transport: transport}
	}

	return tlsClient, err
}

func (this *Client) LoadCert(path string) (err error) {
	if len(path) == 0 {
		return ErrNotFoundCertFile
	}

	cert, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	tlsClient, err := initTLSClient(cert, this.mchId)
	if err != nil {
		return err
	}
	this.tlsClient = tlsClient
	return nil
}

func (this *Client) URLValues(param Param, key string) (value url.Values, err error) {
	var p = param.Params()
	p.Set("appid", this.appId)
	p.Set("mch_id", this.mchId)
	p.Set("nonce_str", GetNonceStr())

	if _, ok := p["notify_url"]; ok == false {
		if len(this.NotifyURL) > 0 {
			p.Set("notify_url", this.NotifyURL)
		}
	}

	var sign = SignMD5(p, key)
	p.Set("sign", sign)
	return p, nil
}

func (this *Client) doRequest(method, url string, param Param, result interface{}) (err error) {
	return this.doRequestWithClient(this.Client, method, url, param, result)
}

func (this *Client) doRequestWithTLS(method, url string, param Param, result interface{}) (err error) {
	if this.tlsClient == nil {
		return ErrNotFoundTLSClient
	}

	return this.doRequestWithClient(this.tlsClient, method, url, param, result)
}

func (this *Client) doRequestWithClient(client *http.Client, method, url string, param Param, result interface{}) (err error) {
	key, err := this.getKey()
	if err != nil {
		return err
	}

	p, err := this.URLValues(param, key)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, strings.NewReader(URLValueToXML(p)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	resp, err := client.Do(req)
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

	if _, err := VerifyResponseData(data, key); err != nil {
		return err
	}

	err = xml.Unmarshal(data, result)

	return err
}

func (this *Client) DoRequest(method, url string, param Param, results interface{}) (err error) {
	return this.doRequest(method, url, param, results)
}

func (this *Client) getKey() (key string, err error) {
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

func (this *Client) SignMD5(param url.Values) (sign string) {
	return SignMD5(param, this.apiKey)
}

func (this *Client) getSignKey(apiKey string) (key string, err error) {
	var p = make(url.Values)
	p.Set("mch_id", this.mchId)
	p.Set("nonce_str", GetNonceStr())

	var sign = SignMD5(p, apiKey)
	p.Set("sign", sign)

	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/sandboxnew/pay/getsignkey", strings.NewReader(URLValueToXML(p)))
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
	var signKey *GetSignKeyRsp
	if err = xml.Unmarshal(data, &signKey); err != nil {
		return "", err
	}
	return signKey.SandboxSignKey, nil
}

func (this *Client) BuildAPI(paths ...string) string {
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

func URLValueToXML(m url.Values) string {
	var xmlBuffer = &bytes.Buffer{}
	xmlBuffer.WriteString("<xml>")

	for key := range m {
		var value = m.Get(key)
		if key == "total_fee" || key == "refund_fee" || key == "execute_time" {
			xmlBuffer.WriteString("<" + key + ">" + value + "</" + key + ">")
		} else {
			xmlBuffer.WriteString("<" + key + "><![CDATA[" + value + "]]></" + key + ">")
		}
	}
	xmlBuffer.WriteString("</xml>")
	return xmlBuffer.String()
}

func SignMD5(param url.Values, key string) (sign string) {
	var pList = make([]string, 0, 0)
	for key := range param {
		var value = param.Get(key)
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
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

func VerifyResponseData(data []byte, key string) (ok bool, err error) {
	var param = make(XMLMap)
	err = xml.Unmarshal(data, &param)
	if err != nil {
		return false, err
	}

	return VerifyResponseValues(url.Values(param), key)
}

func VerifyResponseValues(param url.Values, key string) (bool, error) {
	// 处理错误信息
	var code = param.Get("return_code")
	if code == K_RETURN_CODE_FAIL {
		var msg = param.Get("return_msg")
		if msg == "" {
			msg = param.Get("retmsg")
		}
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

	var sign2 = SignMD5(param, key)
	if sign == sign2 {
		return true, nil
	}
	return false, errors.New("签名验证失败")
}

func GetNonceStr() (nonceStr string) {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 32; i++ {
		idx := r.Intn(len(chars) - 1)
		nonceStr += chars[idx : idx+1]
	}
	return nonceStr
}

func pkcs12ToPem(p12 []byte, password string) (cert tls.Certificate, err error) {
	blocks, err := pkcs12.ToPEM([]byte(p12), password)

	if err != nil {
		return cert, err
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err = tls.X509KeyPair(pemData, pemData)
	return cert, err
}
