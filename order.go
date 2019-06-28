package wxpay

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	kUnifiedOrder = "/pay/unifiedorder"
	kOrderQuery   = "/pay/orderquery"
	kCloseOrder   = "/pay/closeorder"
	kDownloadBill = "/pay/downloadbill"
)

// UnifiedOrder https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (this *Client) UnifiedOrder(param UnifiedOrderParam) (result *UnifiedOrderRsp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(kUnifiedOrder), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

// AppPay APP 支付  https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_12&index=2#
func (this *Client) AppPay(param UnifiedOrderParam) (result *PayInfo, err error) {
	param.TradeType = K_TRADE_TYPE_APP
	rsp, err := this.UnifiedOrder(param)
	if err != nil {
		return nil, err
	}

	if rsp != nil {
		result = &PayInfo{}
		result.AppId = param.AppId
		result.PartnerId = this.mchId
		result.PrepayId = rsp.PrepayId
		result.Package = "Sign=WXPay"
		result.NonceStr = GetNonceStr()
		result.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
		result.SignType = kSignTypeMD5

		var p = url.Values{}
		p.Set("appid", result.AppId)
		p.Set("noncestr", result.NonceStr)
		p.Set("partnerid", result.PartnerId)
		p.Set("prepayid", result.PrepayId)
		p.Set("package", result.Package)
		p.Set("timestamp", result.TimeStamp)

		result.Sign = SignMD5(p, this.apiKey)
		result.RawRsp = rsp
	}
	return result, err
}

// JSAPIPay 微信内H5调起支付-公众号支付 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
func (this *Client) JSAPIPay(param UnifiedOrderParam) (result *PayInfo, err error) {
	param.TradeType = K_TRADE_TYPE_JSAPI
	rsp, err := this.UnifiedOrder(param)
	if err != nil {
		return nil, err
	}

	if rsp != nil {
		result = &PayInfo{}
		result.AppId = param.AppId
		result.PartnerId = this.mchId
		result.PrepayId = rsp.PrepayId
		result.Package = fmt.Sprintf("prepay_id=%s", rsp.PrepayId)
		result.NonceStr = GetNonceStr()
		result.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
		result.SignType = kSignTypeMD5

		var p = url.Values{}
		p.Add("appId", result.AppId)
		p.Add("nonceStr", result.NonceStr)
		p.Add("package", result.Package)
		p.Add("signType", result.SignType)
		p.Add("timeStamp", result.TimeStamp)

		result.Sign = SignMD5(p, this.apiKey)
		result.RawRsp = rsp
	}
	return result, err
}

// MiniAppPay 小程序支付 https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=7_7&index=5
func (this *Client) MiniAppPay(param UnifiedOrderParam) (result *PayInfo, err error) {
	return this.JSAPIPay(param)
}

// WebPay H5 支付 https://pay.weixin.qq.com/wiki/doc/api/H5.php?chapter=9_20&index=1
func (this *Client) WebPay(param UnifiedOrderParam) (result *WebPayInfo, err error) {
	param.TradeType = K_TRADE_TYPE_MWEB
	rsp, err := this.UnifiedOrder(param)
	if err != nil {
		return nil, err
	}

	if rsp != nil {
		result = &WebPayInfo{}
		result.MWebURL = rsp.MWebURL
		result.RawRsp = rsp
	}
	return result, err
}

// NativePay NATIVE 扫码支付 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1
func (this *Client) NativePay(param UnifiedOrderParam) (result *NativePayInfo, err error) {
	param.TradeType = K_TRADE_TYPE_NATIVE
	rsp, err := this.UnifiedOrder(param)
	if err != nil {
		return nil, err
	}

	if rsp != nil {
		result = &NativePayInfo{}
		result.CodeURL = rsp.CodeURL
		result.RawRsp = rsp
	}
	return result, err
}

// OrderQuery https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (this *Client) OrderQuery(param OrderQueryParam) (result *OrderQueryRsp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(kOrderQuery), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

// CloseOrder https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
func (this *Client) CloseOrder(param CloseOrderParam) (result *CloseOrderRsp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(kCloseOrder), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

var (
	kXML = []byte("<xml>")
)

// DownloadBill https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_6
func (this *Client) DownloadBill(param DownloadBillParam) (result *DownloadBillRsp, err error) {
	key, err := this.getKey()
	if err != nil {
		return nil, err
	}

	p, err := this.URLValues(param, key)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", this.BuildAPI(kDownloadBill), strings.NewReader(URLValueToXML(p)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	resp, err := this.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if bytes.Index(data, kXML) == 0 {
		err = xml.Unmarshal(data, &result)
	} else {
		if this.isProduction {
			var r = bytes.NewReader(data)
			gr, err := gzip.NewReader(r)
			if err != nil {
				return nil, err
			}
			defer gr.Close()

			if data, err = ioutil.ReadAll(gr); err != nil {
				return nil, err
			}
		}

		result = &DownloadBillRsp{}
		result.ReturnCode = K_RETURN_CODE_SUCCESS
		result.ReturnMsg = "ok"
		result.Data = data
	}

	return result, err
}
