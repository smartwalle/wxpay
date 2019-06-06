/**
* Author: zhanggaoyuancn@163.com
* Date: 2019-05-29
* Time: 23:40
* Software: GoLand
 */

package wxpay

//公众号签约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_1&index=1
func (client *Client) EntrustWeb(param EntrustWebParam) (result *EntrustWebResponse, err error) {
	if err = client.doRequest("POST", client.BuildAPI(kEntrustWeb), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

//h5签约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_16&index=4
func (client *Client) H5EntrustWeb(param H5EntrustWebParam) (result *H5EntrustWebRsponse, err error) {
	if err = client.doRequest("GET", client.BuildAPI(kH5EntrustWeb), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

//支付签约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_13&index=5
func (client *Client) ContratOrder(param ContratOrderParam) (result *ContratOrderResponse, err error) {
	if err = client.doRequest("POST", client.BuildAPI(kContratOrder), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

//签约申请扣款
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_3&index=8
func (client *Client) PapPayApply(param PapPayApplyParam) (result *PapPayApplyResponse, err error) {
	if err = client.doRequest("POST", client.BuildAPI(kPapPayApply), param, &result); err != nil {
		return nil, err
	}
	return result, err
}

//签约申请解约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_4&index=9
func (client *Client) DeleteContract(param DeleteContractParam) (result *DeleteContractResponse, err error) {
	if err = client.doRequest("POST", client.BuildAPI(kDeleteContract), param, &result); err != nil {
		return nil, err
	}
	return result, err
}
