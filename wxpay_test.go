package wxpay

import (
	"os"
	"testing"
)

var client *WXPay

func TestMain(m *testing.M) {
	//client = New("wx7224e0425e3b8654", "hejzqb1aajazfme7v2zs8e349t8f135e", "1491203582", true) // +1
	client = New("wx20fa044851046bbf", "1v4h5g4s8u1x25tf451d025e10geagf2", "1299730801", false)
	os.Exit(m.Run())
}
