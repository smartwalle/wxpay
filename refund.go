package wxpay

const (
	kRefund        = "/secapi/pay/refund"
	kRefundSandbox = "/pay/refund"
)

// Refund https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_4&index=6
func (this *Client) Refund(param RefundParam) (result *RefundRsp, err error) {
	var api = kRefundSandbox
	if this.isProduction {
		api = kRefund
	}
	if err = this.doRequestWithTLS("POST", this.BuildAPI(api), param, &result); err != nil {
		return nil, err
	}
	return result, err
}
