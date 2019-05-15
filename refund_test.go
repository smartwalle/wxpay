package wxpay

import "testing"

func TestWXPay_Refund(t *testing.T) {
	t.Log("========== Refund ==========")
	var p = RefundParam{}
	p.OutTradeNo = "test-11111112"
	p.OutRefundNo = "ddddd"
	p.TotalFee = 10000
	p.RefundFee = 200

	result, err := client.Refund(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.OutTradeNo, result.TransactionId)
}
