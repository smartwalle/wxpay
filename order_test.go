package wxpay

import (
	"testing"
	"fmt"
)

func TestWXPay_UnifiedOrder(t *testing.T) {
	fmt.Println("========== UnifiedOrder ==========")
	var p = &UnifiedOrderParam{}
	p.Body = "test product"
	p.NotifyURL = "http://www.test.com"
	p.TradeType = "APP"
	p.SpbillCreateIP = "220.112.233.229"
	p.TotalFee = 10
	p.OutTradeNo = "test-111"

	p.StoreInfo = &StoreInfo{}
	p.StoreInfo.Id = "testidtestid"

	result, err := client.UnifiedOrder(p)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(result.PrepayId)
}