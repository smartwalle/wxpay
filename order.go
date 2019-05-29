package wxpay

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
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

	req, err := http.NewRequest("POST", this.BuildAPI(kDownloadBill), strings.NewReader(UrlValueToXML(p)))
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
