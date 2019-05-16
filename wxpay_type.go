package wxpay

import (
	"encoding/xml"
	"errors"
	"io"
	"net/url"
)

const (
	kSandboxURL    = "https://api.mch.weixin.qq.com/sandboxnew"
	kProductionURL = "https://api.mch.weixin.qq.com"
)

const (
	K_RETURN_CODE_FAIL    = "FAIL"
	K_RETURN_CODE_SUCCESS = "SUCCESS"
)

const (
	kSignTypeMD5 = "MD5"
)

var (
	ErrNotFoundCertFile  = errors.New("wxpay: not found cert file")
	ErrNotFoundTLSClient = errors.New("wxpay: not found tls client")
)

type Param interface {
	// 返回参数列表
	Params() url.Values
}

type XMLMap url.Values

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m XMLMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		var e xmlMapEntry
		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(m)[e.XMLName.Local] = []string{e.Value}
	}
	return nil
}

func (v XMLMap) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (v XMLMap) Set(key, value string) {
	v[key] = []string{value}
}

func (v XMLMap) Add(key, value string) {
	v[key] = append(v[key], value)
}

func (v XMLMap) Del(key string) {
	delete(v, key)
}

type GetSignKeyParam struct {
	MchId string
}

func (this *GetSignKeyParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("mch_id", this.MchId)
	return m
}

type GetSignKeyRsp struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	SandboxSignKey string `xml:"sandbox_signkey"`
}
