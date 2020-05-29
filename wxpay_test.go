package wxpay

import (
	"fmt"
	"os"
	"testing"
)

var client *Client

func TestMain(m *testing.M) {
	//client = New("wx7224e0425e3b8654", "hejzqb1aajazfme7v2zs8e349t8f135e", "1491203582", true) // +1
	client = New("wx6149efb9af013077", "dc67305f3c154aa698e802357798d8af", "1533582581", true)

	// 加载退款需要的证书
	fmt.Println(client.LoadCert("./tp.p12"))
	os.Exit(m.Run())
}
