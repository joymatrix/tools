package utils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	alipayV3 "github.com/go-pay/gopay/alipay/v3"
	"github.com/go-pay/gopay/pkg/js"
	"github.com/go-pay/xlog"
)

var aliClient *alipayV3.ClientV3
var privateKey = ``
var appPublicCert = ``

var alipayRootCert = ``

var alipayPublicCert = ``
var alipayAppId = ""

func init() {
	var err error
	//privateKey = cert.PrivateKey
	aliClient, err = alipayV3.NewClientV3(alipayAppId, privateKey, true)
	if err != nil {
		fmt.Printf("alipay new client err:%+v\n", err.Error())
		return
	}
	err = aliClient.SetCert([]byte(appPublicCert), []byte(alipayRootCert), []byte(alipayPublicCert))
	if err != nil {
		fmt.Printf("set cert error:%+v\n", err.Error())
		return
	}

}

type AlipayPay struct {
}

func (ap *AlipayPay) TradePay(c *gin.Context, amount float64, tradeNo string) (string, error) {
	strAmount := strconv.FormatFloat(amount, 'f', 2, 64)
	GetLog().Infof("alipay tradeNo:%s, amount:%s", tradeNo, strAmount)
	// 请求参数
	bm := make(gopay.BodyMap)
	bm.Set("subject", "支付宝订单").
		Set("out_trade_no", tradeNo).
		Set("total_amount", strAmount)

	// 创建订单
	aliRsp, err := aliClient.TradePrecreate(c, bm)
	if err != nil {
		GetLog().Errorf("client.TradePrecreate(), err:%v", err)
		return "", err
	}

	xlog.Warnf("aliRsp:%s", js.Marshal(aliRsp))
	if aliRsp.StatusCode != 200 {
		xlog.Errorf("aliRsp.StatusCode:%d", aliRsp.StatusCode)
		return "", err
	}

	GetLog().Infof("aliRsp:%s", js.Marshal(aliRsp))

	GetLog().Infof("aliRsp.QrCode:", aliRsp.QrCode)
	codeUrl := aliRsp.QrCode

	return codeUrl, nil
}

func (ap *AlipayPay) TradeQuery(c *gin.Context, tradeNo string) (string, error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", tradeNo)
	alipayRsp, err := aliClient.TradeQuery(c, bm)
	if err != nil {
		return "", err
	}
	GetLog().Infof("aliRsp:%s", js.Marshal(alipayRsp))

	if alipayRsp.StatusCode == 10003 || alipayRsp.StatusCode == 20000 {
		GetLog().Infof("alipayRsp.TradeStatus:%+v", alipayRsp.TradeStatus)
		return alipayRsp.TradeStatus, nil
	}

	if alipayRsp.StatusCode != 10000 {
		GetLog().Infof("trade query wxRsp.Code err:%+v", alipayRsp.StatusCode)
		return "", err
	}

	// WAIT_BUYER_PAY（交易创建，等待买家付款）
	// TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）
	// TRADE_SUCCESS（交易支付成功）
	// TRADE_FINISHED（交易结束，不可退款）
	GetLog().Infof("alipayRsp.TradeStatus:%+v", alipayRsp.TradeStatus)
	return alipayRsp.TradeStatus, nil
}

func (ap *AlipayPay) TradeClose(c *gin.Context, tradeNo string) error {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", tradeNo)
	alipayRsp, err := aliClient.TradeClose(c, bm)
	if err != nil {
		return err
	}
	if alipayRsp.StatusCode != 10000 {
		GetLog().Infof("trade close alipayRsp.StatusCode err:%+v", alipayRsp.StatusCode)
		err = errors.New("trade close error")
		return err
	}

	return nil
}

func (ap *AlipayPay) TradeNotify(c *gin.Context) error {
	notifyReq, err := alipay.ParseNotifyToBodyMap(c.Request)
	if err != nil {
		GetLog().Error(err)
		return err
	}

	// 支付宝异步通知验签（公钥证书模式）
	ok, err := alipay.VerifySignWithCert(alipayPublicCert, notifyReq)
	if err != nil {
		return err
	}
	fmt.Println("alipay notify ok:", ok)
	// 如果需要，可将 BodyMap 内数据，Unmarshal 到指定结构体指针 ptr
	//err = notifyReq.Unmarshal(ptr)

	return nil
}

func GetAlipayPay() *AlipayPay {
	return &AlipayPay{}
}
