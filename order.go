package wxpay

import "errors"

const (
	k_UNIFIED_ORDER_URL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

// UnifiedOrder https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (this *WXPay) UnifiedOrder(param *UnifiedOrderParam) (results *UnifiedOrderResp, err error) {
	if err = this.doRequest("POST", k_UNIFIED_ORDER_URL, param, &results); err != nil {
		return nil, err
	}
	if results.ReturnCode == K_RETURN_CODE_FAIL {
		return nil, errors.New(results.ReturnMsg)
	}
	return results, err
}
