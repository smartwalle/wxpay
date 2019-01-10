package wxpay

const (
	kUnifiedOrder = "/pay/unifiedorder"
	kOrderQuery   = "/pay/orderquery"
	kCloseOrder   = "/pay/closeorder"
)

// UnifiedOrder https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (this *WXPay) UnifiedOrder(param UnifiedOrderParam) (result *UnifiedOrderResp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(kUnifiedOrder), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

// OrderQuery https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (this *WXPay) OrderQuery(param OrderQueryParam) (result *OrderQueryResp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(kOrderQuery), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

// CloseOrder https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
func (this *WXPay) CloseOrder(param CloseOrderParam) (result *CloseOrderResp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(kCloseOrder), param, &result); err != nil {
		return nil, err
	}
	return result, err
}
