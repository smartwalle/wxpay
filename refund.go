package wxpay

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

const (
	kRefund        = "/secapi/pay/refund"
	kRefundSandbox = "/pay/refund"
	kRefundQuery   = "/pay/refundquery"
)

// Refund https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_4&index=6
func (this *Client) Refund(param RefundParam) (result *RefundRsp, err error) {
	var api = kRefundSandbox
	if this.isProduction {
		api = kRefund
	}
	if err = this.doRequestWithTLS("POST", this.BuildAPI(api), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

// RefundQuery https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_5&index=7
func (this *Client) RefundQuery(param RefundQueryParam) (result *RefundQueryRsp, err error) {
	body, err := this.doRequestWithClient(this.Client, "POST", this.BuildAPI(kRefundQuery), param, &result)
	if err != nil {
		return nil, err
	}

	if result != nil {
		var infoMap = make(XMLMap)
		xml.Unmarshal(body, &infoMap)

		for i := 0; i < result.RefundCount; i++ {
			var info = &RefundInfo{}
			info.OutRefundNo = infoMap.Get(fmt.Sprintf("out_refund_no_%d", i))
			info.RefundAccount = infoMap.Get(fmt.Sprintf("refund_account_%d", i))
			info.RefundChannel = infoMap.Get(fmt.Sprintf("refund_channel_%d", i))
			info.RefundFee, _ = strconv.Atoi(infoMap.Get(fmt.Sprintf("refund_fee_%d", i)))
			info.RefundId = infoMap.Get(fmt.Sprintf("refund_id_%d", i))
			info.RefundRecvAccount = infoMap.Get(fmt.Sprintf("refund_recv_accout_%d", i))
			info.RefundStatus = infoMap.Get(fmt.Sprintf("refund_status_%d", i))
			info.RefundSuccessTime = infoMap.Get(fmt.Sprintf("refund_success_time_%d", i))
			result.RefundInfos = append(result.RefundInfos, info)
		}
	}
	return result, err
}
