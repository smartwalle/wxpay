/**
* Author: zhanggaoyuancn@163.com
* Date: 2019-05-29
* Time: 23:49
* Software: GoLand
 */

package wxpay

import (
	"fmt"
	"net/url"
)

const (
	kEntrustWeb     = "/papay/entrustweb"     //公众号申请签约
	kContratOrder   = "/pay/contractorder"    //支付中签约
	kPapPayApply    = "/pay/pappayapply"      //签约申请扣款
	kDeleteContract = "/papay/deletecontract" //申请解约
)

//公众号签约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_1&index=1
type EntrustWebParam struct {
	PlanId                 string //是 协议模板id，设置路径见开发步骤docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=17_3。
	ContractCode           string //是 商户侧的签约协议号，由商户生成
	RequestSerial          int64  //是 商户请求签约时的序列号，要求唯一性。序列号主要用于排序，不作为查询条件，纯数字,范围不能超过Int64的范围（9223372036854775807）。
	ContractDisplayAccount string //是 签约用户的名称，用于页面展示，，参数值不支持UTF8非3字节编码的字符，例如表情符号，所以请勿传微信昵称到该字段
	NotifyUrl              string //是 用于接收签约成功消息的回调通知地址，以http或https开头。请对notify_url参数值进行encode处理,注意是对参数值进行encode，例如参数为notify_url=“https://weixin.qq.com”，则需要encode的内容为https://weixin.qq.com
	//Version                string //是 固定值1.0
	Sign      string //是 签名 docs:https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=4_3
	Timestamp string //是 系统当前时间，10位
	ReturnWeb int    //否 1表示返回签约页面的referrer url, 不填或获取不到referrer则不返回; 跳转referrer url时会自动带上参数from_wxpay=1
}

func (entrustWeb EntrustWebParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("notify_url", entrustWeb.NotifyUrl)
	m.Set("plan_id", entrustWeb.PlanId)
	m.Set("contract_code", entrustWeb.ContractCode)
	m.Set("request_serial", fmt.Sprintf("%d", entrustWeb.RequestSerial))
	m.Set("contract_display_account", entrustWeb.ContractDisplayAccount)
	m.Set("version", "1.0")
	m.Set("timestamp", entrustWeb.Timestamp)
	m.Set("return_web", fmt.Sprintf("%d", entrustWeb.ReturnWeb))
	return m
}

//签约成功后，微信会把相关签约结果异步发送给商户，返回的url为调用上述签约接口时填写的notify_url字段。商户在收到签约结果通知后，需进行接收处理并返回应答
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_1&index=1
type EntrustWebResponse struct {
	ReturnCode              string `xml:"return_code"`
	ReturnMsg               string `xml:"return_msg"`
	ResultCode              string `xml:"result_code"`
	MchId                   string `xml:"mch_id"`
	ContractCode            string `xml:"contract_code"`
	PlanId                  string `xml:"plan_id"`
	OpenId                  string `xml:"open_id"`
	Sign                    string `xml:"sign"`
	ChangeType              string `xml:"change_type"`
	OperateTime             string `xml:"operate_time"`
	ContractId              string `xml:"contract_id"`
	ContractExpiredTime     string `xml:"contract_expired_time"`
	ContractTerminationMode int    `xml:"contract_termination_mode"`
	RequestSerial           int64  `xml:"request_serial"`
}

//支付签约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_13&index=5
type ContratOrderParam struct {
	ContractMchId          string //签约商户号，必须与mch_id一致
	ContractAppId          string //签约公众号，必须与appid一致
	OutTradeNo             string //商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	DeviceInfo             string //终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	Body                   string //商品或支付单简要描述
	Detail                 string //商品名称明细列表
	Attach                 string //附加数据,在查询API和支付通知中原样返回,该字段主要用于商户携带订单的自定义数据
	NotifyUrl              string //支付回调通知的url
	TotalFee               int    //订单总金额，单位为分
	SpbillCreateIp         string //APP和网页支付提交用户端ip,Native支付填调用微信支付API的机器IP.
	TimeStart              string //订单生成时间,格式为yyyyMMddHHmmss,如2009年12月25日9点10分10秒表示为20091225091010. 其他详见时间规则
	TimeExpire             string //订单失效时间,格式为yyyyMMddHHmmss,如2009年12月27日9点10分10秒表示为20091227091010. 其他详见时间规则 注意：最短失效时间间隔必须大于5分钟
	GoodsTag               string //商品标记,代金券或立减优惠功能的参数
	TradeType              string //取值如下：JSAPI,NATIVE,APP,MWEB
	ProductId              string //trade_type=NATIVE,此参数必传. 此id为二维码中包含的商品ID,商户自行定义.
	LimitPay               string //no_credit--指定不能使用信用卡支付
	OpenId                 string //trade_type=JSAPI,此参数必传，用户在商户appid下的唯一标识.
	PlanId                 int    //协议模板id
	ContractCode           string //签约协议号
	RequestSerial          int64  //商户请求签约时的序列号，要求唯一性。序列号主要用于排序，不作为查询条件，纯数字,范围不能超过Int64的范围（9223372036854775807）。
	ContractDisplayAccount string //签约用户的名称,用于页面展示，参数值不支持UTF8非3字节编码的字符，例如表情符号，所以请勿传微信昵称到该字段
	ContractNotifyUrl      string //签约信息回调通知的url
}

type ContratOrderResponse struct {
	ReturnCode string //SUCCESS/FAIL 此字段是通信标识,非交易标识,交易是否成功需要查看result_code来判断.
	ReturnMsg  string //返回信息,如非空,为错误原因 /签名失败/参数格式校验错误

	//以下字段在return_code为SUCCESS的时候返回
	ResultCode         string //SUCCESS/FAIL
	AppId              string // appid是商户在微信申请公众号或移动应用成功后分配的帐号ID，登录平台为mp.weixin.qq.com或open.weixin.qq.com
	MchId              string //商户号是商户在微信申请微信支付成功后分配的帐号ID，登录平台为pay.weixin.qq.com
	NonceStr           string //随机字符串,不长于32位.
	Sign               string //签名规则详见签名生成算法 注：所有参数都是encode前做签名.
	ErrCode            string //错误返回的错误代码
	ErrCodeDes         string //错误返回的信息描述
	ContractResultCode string //预签约结果
	ContractErrCode    string //预签约错误代码
	ContractErrCodeDes string //预签约错误描述

	//以下字段在return_code 和result_code都为SUCCESS的时候有返回
	PrepayId               string //微信生成的预支付回话标识,用于后续接口调用中使用,该值有效期为2小时.
	TradeType              string //调用接口提交的交易类型，取值如下：JSAPI,NATIVE,APP
	CodeUrl                string //trade_type为NATIVE是有返回,可将该参数值生成二维码展示出来进行扫码支付
	PlanId                 int    // 商户在微信商户平台设置的代扣协议模板id
	RequestSerial          uint64 //商户请求签约时的序列号,商户侧须唯一
	ContractCode           string //商户请求签约时传入的签约协议号,商户侧须唯一
	ContractDisplayAccount string //签约用户的名称,用于页面展示
	MwebUrl                string //mweb_url为拉起微信支付收银台的中间页面，可通过访问该url来拉起微信客户端，完成支付,mweb_url的有效期为5分钟
	OutTradeNo             string //商户订单号
}

func (contratOrder ContratOrderParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("contract_mchid", contratOrder.ContractMchid)
	m.Set("contract_appid", contratOrder.ContractAppid)
	m.Set("out_trade_no", contratOrder.OutTradeNo)
	m.Set("device_info", contratOrder.DeviceInfo)
	m.Set("body", contratOrder.Body)
	m.Set("detail", contratOrder.Detail)
	m.Set("attach", contratOrder.Attach)
	m.Set("notify_url", contratOrder.NotifyUrl)
	m.Set("total_fee", fmt.Sprintf("%d", contratOrder.TotalFee))
	m.Set("spbill_create_ip", contratOrder.SpbillCreateIp)
	m.Set("time_start", contratOrder.TimeStart)
	m.Set("time_expire", contratOrder.TimeExpire)
	m.Set("goods_tag", contratOrder.GoodsTag)
	m.Set("trade_type", contratOrder.TradeType)
	m.Set("product_id", contratOrder.ProductId)
	m.Set("limit_pay", contratOrder.LimitPay)
	m.Set("openid", contratOrder.OpenId)
	m.Set("plan_id", fmt.Sprintf("%d", contratOrder.PlanId))
	m.Set("contract_code", contratOrder.ContractCode)
	m.Set("request_serial", fmt.Sprintf("%d", contratOrder.RequestSerial))
	m.Set("contract_display_account", contratOrder.ContractDisplayAccount)
	m.Set("contract_notify_url", contratOrder.ContractNotifyUrl)
	return m
}

//签约申请扣款
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_3&index=8
type PapPayApplyParam struct {
	Body           string //是 商品或支付单简要描述
	Detail         string //否 商品名称明细列表
	Attach         string //否 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	OutTradeNo     string //是 商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	TotalFee       int    //是 订单总金额，单位为分，只能为整数，详见支付金额
	FeeType        string //否 符合ISO 4217标准的三位字母代码，默认人民币：CNY
	SpbillCreateIp string //是 调用微信支付API的机器IP
	GoodsTag       string //是 商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
	NotifyUrl      string //是 接受扣款结果异步回调通知的url
	TradeType      string //是 交易类型PAP-微信委托代扣支付
	ContractId     string //是 签约成功后，微信返回的委托代扣协议id
	Receipt        string //否 Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
}

//请求参数
//扣款接口请求成功，返回success仅代表扣款申请受理成功，不代表扣款成功。扣款是否成功以支付通知的结果为准。
func (papPayApply PapPayApplyParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("body", papPayApply.Body)
	m.Set("detail", papPayApply.Detail)
	m.Set("attach", papPayApply.Attach)
	m.Set("out_trade_no", papPayApply.OutTradeNo)
	m.Set("total_fee", fmt.Sprintf("%d", papPayApply.TotalFee))
	m.Set("fee_type", papPayApply.FeeType)
	m.Set("spbill_create_ip", papPayApply.SpbillCreateIp)
	m.Set("goods_tag", papPayApply.GoodsTag)
	m.Set("notify_url", papPayApply.NotifyUrl)
	m.Set("trade_type", papPayApply.TradeType)
	m.Set("contract_id", papPayApply.ContractId)
	m.Set("receipt", papPayApply.Receipt)
	return m
}

//response
//err_code: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_3&index=8
type PapPayApplyResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

//请求参数
//签约申请解约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_4&index=9
type DeleteContractParam struct {
	PlanId                    string //否 商户在微信商户平台配置的代扣模版id，选择plan_id+contract_code解约，则此参数必填
	ContractId                string //否 委托代扣签约成功后由微信返回的委托代扣协议id，选择contract_id解约，则此参数必填
	ContractCode              string //否 商户请求签约时传入的签约协议号，商户侧须唯一。选择plan_id+contract_code解约，则此参数必填
	ContractTerminationRemark string //是 解约原因的备注说明，如：签约信息有误，须重新签约
	//version //固定值1.0
}

func (deleteContract DeleteContractParam) Params() url.Values {
	var m = make(url.Values)
	m.Set("plan_id", deleteContract.PlanId)
	m.Set("contract_id", deleteContract.ContractId)
	m.Set("contract_code", deleteContract.ContractCode)
	m.Set("contract_termination_remark", deleteContract.ContractTerminationRemark)
	m.Set("version", "1.0")
	return m
}

//返回参数
//签约申请解约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_4&index=9
type DeleteContractResponse struct {
	ReturnCode   string `xml:"return_code"`   //SUCCESS/FAIL 此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
	ReturnMsg    string `xml:"return_msg"`    //返回信息，如非空，为错误原因 签名失败 参数格式校验错误
	AppId        string `xml:"appid"`         //微信支付分配的公众账号id
	MchId        string `xml:"mch_id"`        //微信支付分配的商户号
	ContractId   string `xml:"contract_id"`   //委托代扣签约成功后由微信返回的委托代扣协议id
	PlanId       string `xml:"plan_id"`       //商户在微信商户平台设置的代扣协议模板id
	ContractCode string `xml:"contract_code"` //商户请求签约时传入的签约协议号，商户侧须唯一
	ResultCode   string `xml:"result_code"`   //SUCCESS/FAIL
	ErrCode      string `xml:"err_code"`      //错误码
	ErrCodeDes   string `xml:"err_code_des"`  //错误码描述
	Sign         string `xml:"sign"`          //签名
}
