package wxpay

import (
	"fmt"
	"net/url"
)

////////////////////////////////////////////////////////////////////////////////
// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_4&index=6
type RefundParam struct {
	NotifyURL     string // 否 异步接收微信支付退款结果通知的回调地址，通知URL必须为外网可访问的url，不允许带参数, 如果参数中传了notify_url，则商户平台上配置的回调地址将不会生效。
	SignType      string // 否 签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	TransactionId string // 微信生成的订单号，在支付通知中有返回
	OutTradeNo    string // 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	OutRefundNo   string // 是 商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	TotalFee      int    // 是 订单总金额，单位为分，只能为整数，详见支付金额
	RefundFee     int    // 是 退款总金额，订单总金额，单位为分，只能为整数，详见支付金额
	RefundFeeType string // 否 退款货币类型，需与支付一致，或者不填。符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	RefundDesc    string // 否 若商户传入，会在下发给用户的退款消息中体现退款原因, 注意：若订单退款金额≤1元，且属于部分退款，则不会在退款消息中体现退款原因
	RefundAccount string // 否 仅针对老资金流商户使用: REFUND_SOURCE_UNSETTLED_FUNDS---未结算资金退款（默认使用未结算资金退款）, REFUND_SOURCE_RECHARGE_FUNDS---可用余额退款
}

func (this RefundParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("notify_url", this.NotifyURL)
	if len(this.SignType) == 0 {
		this.SignType = kSignTypeMD5
	}
	m.Set("sign_type", this.SignType)
	if len(this.TransactionId) > 0 {
		m.Set("transaction_id", this.TransactionId)
	}
	if len(this.OutTradeNo) > 0 {
		m.Set("out_trade_no", this.OutTradeNo)
	}
	m.Set("out_refund_no", this.OutRefundNo)
	m.Set("total_fee", fmt.Sprintf("%d", this.TotalFee))
	m.Set("refund_fee", fmt.Sprintf("%d", this.RefundFee))
	if len(this.RefundFeeType) > 0 {
		m.Set("refund_fee_type", this.RefundFeeType)
	}
	if len(this.RefundDesc) > 0 {
		m.Set("refund_desc", this.RefundDesc)
	}
	if len(this.RefundAccount) > 0 {
		m.Set("refund_account", this.RefundAccount)
	}
	return m
}

type RefundRsp struct {
	ReturnCode          string `xml:"return_code"`
	ReturnMsg           string `xml:"return_msg"`
	ResultCode          string `xml:"result_code"`
	ErrCode             string `xml:"err_code"`
	ErrCodeDes          string `xml:"err_code_des"`
	AppId               string `xml:"appid"`
	MCHId               string `xml:"mch_id"`
	NonceStr            string `xml:"nonce_str"`
	Sign                string `xml:"sign"`
	TransactionId       string `xml:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no"`
	OutRefundNo         string `xml:"out_refund_no"`
	RefundId            string `xml:"refund_id"`
	RefundFee           int    `xml:"refund_fee"`
	SettlementRefundFee int    `xml:"settlement_refund_fee"`
	TotalFee            int    `xml:"total_fee"`
	SettlementTotalFee  int    `xml:"settlement_total_fee"`
	FeeType             string `xml:"fee_type"`
	CashFee             int    `xml:"cash_fee"`
	CashFeeType         string `xml:"cash_fee_type"`
	CashRefundFee       int    `xml:"cash_refund_fee"`
	CouponRefundFee     int    `xml:"coupon_refund_fee"`
	CouponRefundCount   int    `xml:"coupon_refund_count"`
}
