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
func (client *Client) H5EntrustWeb() {

}

//支付签约
//docs: https://pay.weixin.qq.com/wiki/doc/api/pap.php?chapter=18_13&index=5
func (client *Client) ContratOrder() {

}
