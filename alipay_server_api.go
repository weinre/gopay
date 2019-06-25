//==================================
//  * Name：Jerry
//  * DateTime：2019/6/18 19:24
//  * Desc：
//==================================
package gopay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//解析支付宝支付完成后的Notify信息
func ParseAliPayNotifyResult(req *http.Request) (notifyRsp *AliPayNotifyRequest, err error) {
	notifyRsp = new(AliPayNotifyRequest)
	defer req.Body.Close()
	err = json.NewDecoder(req.Body).Decode(notifyRsp)
	if err != nil {
		return nil, err
	}
	return
}

//支付通知的签名验证和参数签名后的Sign
//    alipayPublickKey：支付宝公钥
//    notifyRsp：利用 gopay.ParseAliPayNotifyResult() 得到的结构体
//    返回参数ok：是否验证通过
//    返回参数sign：根据参数计算的sign值，非支付宝返回参数中的Sign
func VerifyAliPayResultSign(alipayPublickKey string, notifyRsp *AliPayNotifyRequest) (ok bool, sign string) {
	body := make(BodyMap)
	body.Set("notify_time", notifyRsp.NotifyTime)
	body.Set("notify_type", notifyRsp.NotifyType)
	body.Set("notify_id", notifyRsp.NotifyId)
	body.Set("app_id", notifyRsp.AppId)
	body.Set("charset", notifyRsp.Charset)
	body.Set("version", notifyRsp.Version)
	body.Set("trade_no", notifyRsp.TradeNo)
	body.Set("out_trade_no", notifyRsp.OutTradeNo)
	body.Set("out_biz_no", notifyRsp.OutBizNo)
	body.Set("buyer_id", notifyRsp.BuyerId)
	body.Set("buyer_logon_id", notifyRsp.BuyerLogonId)
	body.Set("seller_id", notifyRsp.SellerId)
	body.Set("seller_email", notifyRsp.SellerEmail)
	body.Set("trade_status", notifyRsp.TradeStatus)
	body.Set("total_amount", notifyRsp.TotalAmount)
	body.Set("receipt_amount", notifyRsp.ReceiptAmount)
	body.Set("invoice_amount", notifyRsp.InvoiceAmount)
	body.Set("buyer_pay_amount", notifyRsp.BuyerPayAmount)
	body.Set("point_amount", notifyRsp.PointAmount)
	body.Set("refund_fee", notifyRsp.RefundFee)
	body.Set("subject", notifyRsp.Subject)
	body.Set("body", notifyRsp.Body)
	body.Set("gmt_create", notifyRsp.GmtCreate)
	body.Set("gmt_payment", notifyRsp.GmtPayment)
	body.Set("gmt_refund", notifyRsp.GmtRefund)
	body.Set("gmt_close", notifyRsp.GmtClose)
	body.Set("fund_bill_list", jsonToString(notifyRsp.FundBillList))
	body.Set("passback_params", notifyRsp.PassbackParams)
	body.Set("voucher_detail_list", jsonToString(notifyRsp.VoucherDetailList))

	newBody := make(BodyMap)
	for k, v := range body {
		if v != null {
			newBody.Set(k, v)
		}
	}

	sign, err := getRsaSign(newBody, alipayPublickKey)
	if err != nil {
		return false, ""
	}
	ok = sign == notifyRsp.Sign
	return
}

func jsonToString(v interface{}) (str string) {
	bs, err := json.Marshal(v)
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}
	//log.Println("string:", string(bs))
	return string(bs)
}

//格式化秘钥
func FormatPrivateKey(privateKey string) (pKey string) {
	buffer := new(bytes.Buffer)
	buffer.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")

	rawLen := 64
	keyLen := len(privateKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen

	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(privateKey[start:])
		} else {
			buffer.WriteString(privateKey[start:end])
		}
		buffer.WriteString("\n")
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END RSA PRIVATE KEY-----\n")
	pKey = buffer.String()
	return
}

//格式化秘钥
func FormatAliPayPublicKey(publickKey string) (pKey string) {
	buffer := new(bytes.Buffer)
	buffer.WriteString("-----BEGIN PUBLIC KEY-----\n")

	rawLen := 64
	keyLen := len(publickKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen

	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(publickKey[start:])
		} else {
			buffer.WriteString(publickKey[start:end])
		}
		buffer.WriteString("\n")
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END PUBLIC KEY-----\n")
	pKey = buffer.String()
	return
}