package wxpay

import (
	"fmt"
	"testing"
)

func TestWXPay_UnifiedOrder(t *testing.T) {
	fmt.Println("========== UnifiedOrder ==========")
	var p = UnifiedOrderParam{}
	p.Body = "test product"
	p.NotifyURL = "http://www.test.com"
	p.TradeType = K_TRADE_TYPE_NATIVE
	p.SpbillCreateIP = "202.105.107.18"
	p.TotalFee = 101
	p.OutTradeNo = "test-11111112"

	result, err := client.UnifiedOrder(p)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(result.PrepayId, result.CodeURL)
}

func TestWXPay_OrderQuery(t *testing.T) {
	fmt.Println("========== OrderQuery ==========")
	var p = OrderQueryParam{}
	p.OutTradeNo = "test-11111112"

	result, err := client.OrderQuery(p)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(result.TradeState)
}
