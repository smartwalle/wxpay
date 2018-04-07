package wxpay

const (
	k_UNIFIED_ORDER_URL = "/pay/unifiedorder"
	k_ORDER_QUERY       = "/pay/orderquery"
)

// UnifiedOrder https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (this *WXPay) UnifiedOrder(param UnifiedOrderParam) (result *UnifiedOrderResp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(k_UNIFIED_ORDER_URL), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

// OrderQuery https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
func (this *WXPay) OrderQuery(param OrderQueryParam) (result *OrderQueryResp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(k_ORDER_QUERY), param, &result); err != nil {
		return nil, err
	}
	return result, err
}
