package wxpay

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	K_TRADE_TYPE_JSAPI  = "JSAPI"
	K_TRADE_TYPE_NATIVE = "NATIVE"
	K_TRADE_TYPE_APP    = "APP"
	K_TRADE_TYPE_MWEB   = "MWEB"
)

const (
	K_TRADE_STATE_SUCCESS    = "SUCCESS"    //支付成功
	K_TRADE_STATE_REFUND     = "REFUND"     //转入退款
	K_TRADE_STATE_NOTPAY     = "NOTPAY"     //未支付
	K_TRADE_STATE_CLOSED     = "CLOSED"     //已关闭
	K_TRADE_STATE_REVOKED    = "REVOKED"    //已撤销（刷卡支付）
	K_TRADE_STATE_USERPAYING = "USERPAYING" //用户支付中
	K_TRADE_STATE_PAYERROR   = "PAYERROR"   //支付失败(其他原因，如银行返回失败)
)

////////////////////////////////////////////////////////////////////////////////
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
type UnifiedOrderParam struct {
	AppId          string // 是
	NotifyURL      string // 是 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
	Body           string // 是 商品简单描述，该字段请按照规范传递，具体请见参数规定
	OutTradeNo     string // 是 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号
	TotalFee       int    // 是 订单总金额，单位为分，详见支付金额
	SpbillCreateIP string // 是 APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	TradeType      string // 是 取值如下：JSAPI，NATIVE，APP等，说明详见参数规定
	SignType       string // 否 签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	DeviceInfo     string // 否 自定义参数，可以为终端设备号(门店号或收银设备ID)，PC网页或公众号内支付可以传"WEB"
	Detail         string // 否 商品详细描述，对于使用单品优惠的商户，改字段必须按照规范上传，详见“单品优惠参数说明”
	Attach         string // 否 附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	FeeType        string // 否 符合ISO 4217标准的三位字母代码，默认人民币：CNY，详细列表请参见货币类型
	TimeStart      string // 否 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TimeExpire     string // 否 订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则  注意：最短失效时间间隔必须大于5分钟
	GoodsTag       string // 否 订单优惠标记，使用代金券或立减优惠功能时需要的参数，说明详见代金券或立减优惠
	ProductId      string // 否 trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	LimitPay       string // 否 上传此参数no_credit--可限制用户不能使用信用卡支付
	OpenId         string // 否 trade_type=JSAPI时（即公众号支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识。openid如何获取，可参考【获取openid】。企业号请使用【企业号OAuth2.0接口】获取企业号内成员userid，再调用【企业号userid转openid接口】进行转换
	SceneInfo      string // 否
	StoreInfo      *StoreInfo
}

type StoreInfo struct {
	Id       string `json:"id"`        // 门店唯一标识
	Name     string `json:"name"`      // 门店名称
	AreaCode string `json:"area_code"` // 门店所在地行政区划码，详细见《最新县及县以上行政区划代码》
	Address  string `json:"address"`   // 门店详细地址
}

func (this UnifiedOrderParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("appid", this.AppId)
	m.Set("notify_url", this.NotifyURL)
	if len(this.SignType) == 0 {
		this.SignType = kSignTypeMD5
	}
	m.Set("sign_type", this.SignType)
	m.Set("device_info", this.DeviceInfo)
	m.Set("body", this.Body)
	m.Set("detail", this.Detail)
	m.Set("attach", this.Attach)
	m.Set("out_trade_no", this.OutTradeNo)
	m.Set("fee_type", this.FeeType)
	m.Set("total_fee", fmt.Sprintf("%d", this.TotalFee))
	m.Set("spbill_create_ip", this.SpbillCreateIP)
	m.Set("time_start", this.TimeStart)
	m.Set("time_expire", this.TimeExpire)
	m.Set("goods_tag", this.GoodsTag)
	if len(this.TradeType) == 0 {
		this.TradeType = K_TRADE_TYPE_APP
	}
	m.Set("trade_type", this.TradeType)
	m.Set("product_id", this.ProductId)
	m.Set("limit_pay", this.LimitPay)
	m.Set("openid", this.OpenId)

	if this.StoreInfo != nil {
		var storeInfoByte, err = json.Marshal(this.StoreInfo)
		if err == nil {
			this.SceneInfo = "{\"store_info\" :" + string(storeInfoByte) + "}"
			m.Set("scene_info", this.SceneInfo)
		}
	} else if this.SceneInfo != "" {
		m.Set("scene_info", this.SceneInfo)
	}
	return m
}

type UnifiedOrderRsp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MCHId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
	CodeURL    string `xml:"code_url"`
	MWebURL    string `xml:"mweb_url"`
}

// 客户端唤起支付所需要的信息：App 支付、微信内H5调起支付(公众号支付)、小程序支付
// App 支付 - https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_12&index=2
// 微信内H5调起支付 - https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
// 小程序调起支付API - https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=7_7&index=5
type PayInfo struct {
	AppId     string           `json:"app_id"`
	PartnerId string           `json:"partner_id"`
	PrepayId  string           `json:"prepay_id"`
	Package   string           `json:"package"`
	NonceStr  string           `json:"nonce_str"`
	TimeStamp string           `json:"timestamp"`
	Sign      string           `json:"sign"`
	SignType  string           `json:"sign_type"`
	RawRsp    *UnifiedOrderRsp `json:"-"`
}

// WebPayInfo https://pay.weixin.qq.com/wiki/doc/api/H5.php?chapter=9_20&index=1
type WebPayInfo struct {
	MWebURL string           `json:"mweb_url"`
	RawRsp  *UnifiedOrderRsp `json:"-"`
}

// NativePayInfo https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_1
type NativePayInfo struct {
	CodeURL string           `json:"code_url"`
	RawRsp  *UnifiedOrderRsp `json:"-"`
}

////////////////////////////////////////////////////////////////////////////////
// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_2&index=4
type OrderQueryParam struct {
	TransactionId string
	OutTradeNo    string
}

func (this OrderQueryParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("transaction_id", this.TransactionId)
	m.Set("out_trade_no", this.OutTradeNo)
	return m
}

type OrderQueryRsp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MCHId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	OpenId             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	TradeState         string `xml:"trade_state"`
	BankType           string `xml:"bank_type"`
	TotalFee           int    `xml:"total_fee"`
	SettlementTotalFee int    `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            int    `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          int    `xml:"coupon_fee"`
	CouponCount        int    `xml:"coupon_count"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
	TradeStateDesc     string `xml:"trade_state_desc"`
}

////////////////////////////////////////////////////////////////////////////////
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
type CloseOrderParam struct {
	OutTradeNo string // 是 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
}

func (this CloseOrderParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("out_trade_no", this.OutTradeNo)
	return m
}

type CloseOrderRsp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MCHId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

////////////////////////////////////////////////////////////////////////////////
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_6
type DownloadBillParam struct {
	BillDate string `xml:"bill_date"` // 是 下载对账单的日期，格式：20140603
	BillType string `xml:"bill_type"` // 是 ALL，返回当日所有订单信息，默认值；SUCCESS，返回当日成功支付的订单；REFUND，返回当日退款订单；RECHARGE_REFUND，返回当日充值退款订单
	TarType  string `xml:"tar_type"`  // 否 非必传参数，固定值：GZIP，返回格式为.gzip的压缩包账单。不传则默认为数据流形式。
}

func (this DownloadBillParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("bill_date", this.BillDate)
	m.Set("bill_type", this.BillType)
	m.Set("tar_type", this.TarType)
	return m
}

type DownloadBillRsp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Data       []byte `xml:"-"`
}
