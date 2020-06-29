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

func TestWXPay_RefundQuery(t *testing.T) {
	t.Log("========== RefundQuery ==========")
	var p = RefundQueryParam{}
	p.OutTradeNo = "21488302741782528"

	result, err := client.RefundQuery(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.OutTradeNo, result.TransactionId, result.TotalFee, result.RefundFee)
	for _, info := range result.RefundInfos {
		t.Log(info.RefundFee, info.OutRefundNo)
	}

}
