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
	kEntrustWeb = "/papay/entrustweb"
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

func (entrustWeb *EntrustWebParam) Params() url.Values {
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
