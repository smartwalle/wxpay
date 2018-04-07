package wxpay

import (
	"fmt"
	"testing"
)

func TestWXPay_UnifiedOrder(t *testing.T) {
	fmt.Println("========== UnifiedOrder ==========")
	var p = &UnifiedOrderParam{}
	p.Body = "test product"
	p.NotifyURL = "http://www.test.com"
	p.TradeType = K_TRADE_TYPE_APP
	p.SpbillCreateIP = "202.105.107.18"
	p.TotalFee = 101
	p.OutTradeNo = "test-111111"

	result, err := client.UnifiedOrder(p)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(result.PrepayId, result.CodeURL)
}
