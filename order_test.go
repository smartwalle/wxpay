package wxpay

import (
	"testing"
)

func TestWXPay_UnifiedOrder(t *testing.T) {
	t.Log("========== UnifiedOrder ==========")
	var p = UnifiedOrderParam{}
	p.Body = "test product"
	p.NotifyURL = "http://www.test.com"
	p.TradeType = K_TRADE_TYPE_NATIVE
	p.SpbillCreateIP = "202.105.107.18"
	p.TotalFee = 101
	p.OutTradeNo = "test-111111122sdf"

	result, err := client.UnifiedOrder(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.PrepayId, result.CodeURL)
}

func TestWXPay_OrderQuery(t *testing.T) {
	t.Log("========== OrderQuery ==========")
	var p = OrderQueryParam{}
	p.OutTradeNo = "test-11111112"

	result, err := client.OrderQuery(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.TradeState, result.OutTradeNo, result.TransactionId)
}

func TestWXPay_CloseOrder(t *testing.T) {
	t.Log("========== CloseOrder ==========")
	var p = CloseOrderParam{}
	p.OutTradeNo = "test-11111112"

	result, err := client.CloseOrder(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.ReturnCode, result.ReturnMsg)
}
