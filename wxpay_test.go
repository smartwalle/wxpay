package wxpay

import (
	"fmt"
	"os"
	"testing"
)

var client *Client

func TestMain(m *testing.M) {
	client = New("wx143cd4036f7c65c4", "dc67305f3c154aa698e802357798d8af", "1533582581", true)

	// 加载退款需要的证书
	fmt.Println(client.LoadCert("./tp.p12"))
	os.Exit(m.Run())
}
