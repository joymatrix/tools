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
	appId             = "wxac396a3be7a16844"
	mchId             = "1615933939"
	serialNo          = "33DDFEC51BF5412F663B9B56510FD567B625FC68"
	apiv3Key          = "V10987654321younggee12345678910V"
	privateKeyContent = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDiS0mOMU85JaH1
nRsOU57TfKP+3ma7TiEe2jQNHhF1UUW8i1+Hc1NeLe7dh29J5aEa/isu0RG2LHjY
P73dzPkarAdhmW3ETsbBVazUw7RfhbPawR6EjBArRAojqHLXdSTFeDM6LtmwbWD2
JEjiBKyTMtsbnap4ubLhR03huPn8V20Hbq3Zg+U8T4XtnK7jrac2SqkuybKheRX8
OzaefU/nYbxKwyTqKCanAs96ZTFuEceW9Utoa4jb4bqdVVhbc2P5gyyssKfYkfDj
Pdr2X3vg11qq7bhjQ+xLgUiWwXlqnKY0eHxvx0Zdh5PWdh0oYZtbuHrkPPJaf4zu
Zkre9VnXAgMBAAECggEASPlaWP5Ru+4E0n29UdtpZm3VPMVff5tsVtSq4GgH3Ts+
L9UKE1X/Vmmdk9au7recQmYgatKE0ah5t9KmWbZVxmIfZzvhB+MXeRU1zM7nhb5K
B4srWjcIp8sjMeiKUCy4lO10J5kgHiLHl9iPoEM9m6JUwg0QAipwIvGpjdbm1paH
esqCyZ2+fyTypO1mkpPOOB6caGwePE0Ppd/U5woPKPS0H5JnBdKPWBh/XYkBigOl
fMRjJg9BDGPjxI1Rt+F6S359fgfXZJV1m74YVL9HxgfKi+WHdDTtj0094NBqizCg
goeLZbz6aWYbOBmpFHrK8F2i/pYU658KLCtjjpvnQQKBgQD6JJGRqqF1OLVsjJ+E
+8ts8IdC+OBYY9YHgLw5D2aFUrtJre9Fiktj0hupXPm+01kmc+5OupR3O/L8jhM6
+NTe3fmIzNLk72GeZsiWFrZCN0fKjts421wwyhBrqW/qZN66a1kYA7fvr0ToRjL1
jLhTqu0a8BGRyaxbqlebo/5FkQKBgQDnl8MYmJoYpJ36WK3Bd1j/9eczDmRXf4RW
kZjIOeXBN4zHZsvlFtMtcaZAw3Z1KPbM4ncuKwl+pc7oNUfOAU/IOiBEiZR7zSYG
D8g+UoXfZzv5mFQQBSagfY4Q0bigG1Zt/QPMoN5setiG6BJ+bwN7WfXNjRjyaXtZ
ovUUniFU5wKBgQDrvu0mcL6MIG7zp7Brf3bf6+w+lRmylBzRo2VBDZ+chTUXooJ/
cm/M2ubQ/lwtmThLAjWVI0jq+qftl+TNzleo12DmqcsUkfrZc5sVwL/ytfDGGU7I
TgybusQxA1YDfR9gZ+1msZJ3pSJ3GjnKq93IlK2zlo+oa34yQd8hQzRP0QKBgQDQ
5TkbNHq6g7HjoJ2KBocGyd2zVeX4bpMGKuouoNq2v86CBh0gFMiDEyItBKIS59JF
2Hg78qHr1M+e8IBGNzSpnJSCfb6rNM55ZT7vyCvs6QdWCaq5kIvY86dzUFhCQqZh
K3mD2A8Itn4cobQcyzHOz8RBlmXMMo0Ku0xpPoE+PQKBgEFcy2wvv/saKJ+3t2Cw
4ckwskhsQ2CizlAlTCNmBzou+IYfcg3VxS7P0EA/ViUXHTOx1pjhS977FLqJxe08
FU7a3mIzzehJA0VnsZC6lzCqSRxdE4bcEs9Opw+sJS7abJKRU2QeZSJEChqAMUhJ
3oP15BVjOfW79TL/hp8fRZAt
-----END PRIVATE KEY-----`
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
		Set("description", "环卫宝微信支付").
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
