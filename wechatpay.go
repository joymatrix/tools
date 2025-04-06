package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	wechat "github.com/go-pay/gopay/wechat/v3"
)

var (
	appId             = ""
	mchId             = ""
	serialNo          = ""
	apiv3Key          = ""
	privateKeyContent = ``
)
var wechatClient *wechat.ClientV3

func init() {
	var err error
	wechatClient, err = wechat.NewClientV3(mchId, serialNo, apiv3Key, privateKeyContent)
	if err != nil {
		GetLog().Errorf("wechat newclient v3 err:%+v", err.Error())
		return
	}

	// 设置微信平台证书和序列号，如开启自动验签，请忽略此步骤
	//client.SetPlatformCert([]byte(""), "")

	// 启用自动同步返回验签，并定时更新微信平台API证书
	err = wechatClient.AutoVerifySign()
	if err != nil {
		GetLog().Errorf("wechat AutoVerifySign err:%+v", err.Error())
		return
	}
}

type WechatPay struct {
}

func (wp *WechatPay) TradePay(c *gin.Context, amount int, tradeNo string) (string, error) {
	GetLog().Infof("tradeNo:%s", tradeNo)
	expire := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	bm := make(gopay.BodyMap)
	bm.Set("appid", appId).
		Set("description", "xxx").
		Set("out_trade_no", tradeNo).
		Set("time_expire", expire).
		Set("notify_url", "https://www.fmm.ink").
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", amount).
				Set("currency", "CNY")
		})

	wxRsp, err := wechatClient.V3TransactionNative(c, bm)
	if err != nil {
		GetLog().Errorf("V3TransactionNative err:%+v", err.Error())
		return "", err
	}
	if wxRsp.Code != 0 {
		GetLog().Infof("trade pay wxRsp.Code err:%+v", wxRsp.Code)
		err = errors.New("trade pay code is not 0")
		return "", err
	}
	codeUrl := wxRsp.Response.CodeUrl
	GetLog().Infof("trade rsp:%+v", wxRsp.Response.CodeUrl)
	return codeUrl, nil
}

func (wp *WechatPay) TradeQuery(c *gin.Context, tradeNo string) (string, string, error) {
	wxRsp, err := wechatClient.V3TransactionQueryOrder(c, wechat.OutTradeNo, tradeNo)
	if err != nil {
		return "", "", err
	}
	if wxRsp.Code != 0 {
		GetLog().Infof("trade query wxRsp.Code err:%+v", wxRsp.Code)
		return "", "", err
	}
	// SUCCESS：支付成功
	// REFUND：转入退款
	// NOTPAY：未支付
	// CLOSED：已关闭
	// REVOKED：已撤销（仅付款码支付会返回）
	// USERPAYING：用户支付中（仅付款码支付会返回）
	// PAYERROR：支付失败（仅付款码支付会返回）
	GetLog().Infof("wxRsp.Response.TradeState:%+v", wxRsp.Response.TradeState)
	return wxRsp.Response.TradeState, wxRsp.Response.TradeStateDesc, nil
}

func (wp *WechatPay) TradeClose(c *gin.Context, tradeNo string) error {
	wxRsp, err := wechatClient.V3TransactionCloseOrder(c, tradeNo)
	if err != nil {
		return err
	}
	if wxRsp.Code != 0 {
		GetLog().Infof("trade close wxRsp.Code err:%+v", wxRsp.Code)
		err = errors.New("trade close error")
		return err
	}

	return nil
}

func (wp *WechatPay) TradeNotify(c *gin.Context) {
	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		return
	}
	certMap := wechatClient.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		return
	}
	result, err := notifyReq.DecryptPayCipherText(apiv3Key)
	if err != nil {
		return
	}
	GetLog().Infof("notify info:%+v,%+v", result.TradeState, result.TradeStateDesc)
	c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
}

func GetWechatPay() *WechatPay {
	return &WechatPay{}
}
