package wxpay

const (
	kRefund = "/secapi/pay/refund"
)

// Refund https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_4&index=6
func (this *WXPay) Refund(param RefundParam) (result *RefundResp, err error) {
	if err = this.doRequestWithTLS("POST", this.BuildAPI(kRefund), param, &result); err != nil {
		return nil, err
	}
	return result, err
}
