package wxpay

const (
	k_UNIFIED_ORDER_URL = "/pay/unifiedorder"
)

// UnifiedOrder https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func (this *WXPay) UnifiedOrder(param *UnifiedOrderParam) (results *UnifiedOrderResp, err error) {
	if err = this.doRequest("POST", this.BuildAPI(k_UNIFIED_ORDER_URL), param, &results); err != nil {
		return nil, err
	}
	return results, err
}
