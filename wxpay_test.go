package wxpay

import (
	"testing"
	"os"
)

var client *WXPay

func TestMain(m *testing.M) {
	client = New("wx20fa044851046bbf", "1v4h5g4s8u1x25tf451d025e10geagf2", "1299730801")
	os.Exit(m.Run())
}