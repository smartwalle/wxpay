package wxpay

type WXPayParam interface {
	// 返回参数列表
	Params() map[string]interface{}
}
