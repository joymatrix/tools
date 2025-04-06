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
var privateKey = `MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCPKM66YV2cenRT8cVxoTUxuInKGAtoifyyGikso5R67tf9sDwHu9qtXi8LM1sxghI3+a+iDyDjAhAmPurHvjLsHNQUkqvhsjcfjojZEqLtLXCh2WAOjEtEHNZ6iFvSBIs27EdMAf+y0rj6XmJ+z8In1VfXfgOD97kkDiDyGdAiFT6orANaGTpLjCrHWC/pbFC191sCPj4QJMse/hjueUeAJacIr8i8iCuK5IYG1qTEWhW2LE6SHQbjzd3kZTBhurobnpdYQs02WFM9knwFicKYWwDGmzamQYsZS6gygs2s8e4+Zv70XnWxzZmBldsFZreyY9CjNek5Xw6h43cPif2bAgMBAAECggEAcs/3e+K6gNR0lx4/i3IOh2HIoBvIin6f+vagLvzCCBWlg//jJRCzwHbYo9L8QChhFCNbiE05wtXUvdeX07nmfRZhwF4hG1EihFx7xBv+LtlSi3saXpCFjIrUOFFD0ptySwoT5BF2UKRPVfx8YdedjvS7Dkgx6ZSzFwd9xKyPD8VbmtRFuHVGQllhoyACEOnFT2Ef/VM8BxZvsWFs4qmBNPlHmxy5lxspplLApLVJvPi6suikWO11XCPFw6VBXsLNVInzuXo5FB42sXTWnkTpg8JoU6pp9rGMsOv3FqH+OZNBNxHoJ9yuX2qo8WJl07bK4WZ5kO96CNmOLmcKKlHNAQKBgQD7Gx+OmaHBN4MyAlxUYPPX/e6TIqEenqnTvfyh9vSP5g5CrbidorL1AsBdKYCBVDMgLEZ7OloJz2+HbEqdjlos2tHsCfL/924ztz2H2wKy9OY0+82scqcejpHH71ds6LGb3aMOAp76wT0n5fYRehx2rb357P8jWP1NikUTtDDUQwKBgQCR8xeO88TdGm1Tfzp86+31rIwYieAE5ZWB5n6NAEsbq+LiRM/yWWaMezju2S1jfEwYKSudkECKyZbT5noXXXLiZRRTEoynEXIi4Ji2qgr1aWsS6aZaROzvc8XCouaysU6bvN0J1Zd4xtPj81FuudSqGxXFTBtxq+rgdZOrg2yHyQKBgQC2uJjxAlhTKhr8o/0dpWMrE+usA3HsvxXjL4eLMBHsOEK+QH1rr727TiI+aHnLIkMFsVIkT2S/aMPGboWpOrHhm+VPdjnuCtWVKkzK9BJ1uIFfoq+aQd/b+3CXZVFfvb+oJNKG2l9nJoBi7RJuy0W1El6AY/WQeivwZaI4YVF60wKBgHAA9BVacbuUak3nl4UCse0Va7XxKR/Y6HF156xhi0JDGKy4TjUX8qDgd4Kk3DY6z1LDVZtndoCLz9nyR8PijSW8mGpgE9yqgMLPRaL4v8wyCF/NO0KPHp1sZVnHFfAQLdlKiP7hEYs4WqfWtKmapt+cRYrRq0YCvw9ea+L1yrYxAoGBAPiVow4wdwBqNkxB+CmMn38IGW4lmvE84T9LZDE4of0br42bJbaqlRO16QlGO+nHOs9UOW02uxGBLg3SAzKlPiKW1pFT+cckhQJLNW5lKXV+zzKzSiK4AXgnhCvGG2UYXAPOP1m9PAMj3H6uyZd3+mDSgZmI15dZpfB6nh1BvWrk`
var appPublicCert = `-----BEGIN CERTIFICATE-----
MIIEoDCCA4igAwIBAgIQICUDEmQuGQoD0Gb3Hzul5jANBgkqhkiG9w0BAQsFADCBgjELMAkGA1UE
BhMCQ04xFjAUBgNVBAoMDUFudCBGaW5hbmNpYWwxIDAeBgNVBAsMF0NlcnRpZmljYXRpb24gQXV0
aG9yaXR5MTkwNwYDVQQDDDBBbnQgRmluYW5jaWFsIENlcnRpZmljYXRpb24gQXV0aG9yaXR5IENs
YXNzIDEgUjEwHhcNMjUwMzEyMDY1MDEyWhcNMzAwMzExMDY1MDEyWjBoMQswCQYDVQQGEwJDTjEt
MCsGA1UECgwk5YyX5Lqs5pm65a2Q5pe256m656eR5oqA5pyJ6ZmQ5YWs5Y+4MQ8wDQYDVQQLDAZB
bGlwYXkxGTAXBgNVBAMMEDIwODg5MzE5MDQ3MjY2MTAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQCPKM66YV2cenRT8cVxoTUxuInKGAtoifyyGikso5R67tf9sDwHu9qtXi8LM1sxghI3
+a+iDyDjAhAmPurHvjLsHNQUkqvhsjcfjojZEqLtLXCh2WAOjEtEHNZ6iFvSBIs27EdMAf+y0rj6
XmJ+z8In1VfXfgOD97kkDiDyGdAiFT6orANaGTpLjCrHWC/pbFC191sCPj4QJMse/hjueUeAJacI
r8i8iCuK5IYG1qTEWhW2LE6SHQbjzd3kZTBhurobnpdYQs02WFM9knwFicKYWwDGmzamQYsZS6gy
gs2s8e4+Zv70XnWxzZmBldsFZreyY9CjNek5Xw6h43cPif2bAgMBAAGjggEpMIIBJTAfBgNVHSME
GDAWgBRxB+IEYRbk5fJl6zEPyeD0PJrVkTAdBgNVHQ4EFgQUDuJi1q4DZCQmvQbl2AyuLpyG0G4w
QAYDVR0gBDkwNzA1BgdggRwBbgEBMCowKAYIKwYBBQUHAgEWHGh0dHA6Ly9jYS5hbGlwYXkuY29t
L2Nwcy5wZGYwDgYDVR0PAQH/BAQDAgbAMC8GA1UdHwQoMCYwJKAioCCGHmh0dHA6Ly9jYS5hbGlw
YXkuY29tL2NybDk5LmNybDBgBggrBgEFBQcBAQRUMFIwKAYIKwYBBQUHMAKGHGh0dHA6Ly9jYS5h
bGlwYXkuY29tL2NhNi5jZXIwJgYIKwYBBQUHMAGGGmh0dHA6Ly9jYS5hbGlwYXkuY29tOjgzNDAv
MA0GCSqGSIb3DQEBCwUAA4IBAQBQ8/8Xc4BlDvfuOVg/wCnQ4hbAAxKMvPKBPkESWFpRqYXZ6uki
ugr2b7tz4Cbbqpee69qwIhqrllbHzSdXRLqiUre3QqIEoU5SFrwZfm9Kt4eREta4xv0Q9Woj6epM
OM4aBIbYakKz+u9zQgY9B8XZTQvHrSxdbEe0KfcnkSjRs+TLVkxPL+qExKr/4OWOSGXcws4opEZN
Hoqett5QkxNTI6axewCLJjmXEPN7Vul18/L7VIyLfGPYIzvhZdZuD6qeXTHoRx5RWu2IWfr4qz8H
RDFu5D3RSUmmWZArXzTahyiWpS6NvsDHTHWAyhG6JQ9EjEfXr6gus2NdA33HcHLQ
-----END CERTIFICATE-----`

var alipayRootCert = `-----BEGIN CERTIFICATE-----
MIIBszCCAVegAwIBAgIIaeL+wBcKxnswDAYIKoEcz1UBg3UFADAuMQswCQYDVQQG
EwJDTjEOMAwGA1UECgwFTlJDQUMxDzANBgNVBAMMBlJPT1RDQTAeFw0xMjA3MTQw
MzExNTlaFw00MjA3MDcwMzExNTlaMC4xCzAJBgNVBAYTAkNOMQ4wDAYDVQQKDAVO
UkNBQzEPMA0GA1UEAwwGUk9PVENBMFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAE
MPCca6pmgcchsTf2UnBeL9rtp4nw+itk1Kzrmbnqo05lUwkwlWK+4OIrtFdAqnRT
V7Q9v1htkv42TsIutzd126NdMFswHwYDVR0jBBgwFoAUTDKxl9kzG8SmBcHG5Yti
W/CXdlgwDAYDVR0TBAUwAwEB/zALBgNVHQ8EBAMCAQYwHQYDVR0OBBYEFEwysZfZ
MxvEpgXBxuWLYlvwl3ZYMAwGCCqBHM9VAYN1BQADSAAwRQIgG1bSLeOXp3oB8H7b
53W+CKOPl2PknmWEq/lMhtn25HkCIQDaHDgWxWFtnCrBjH16/W3Ezn7/U/Vjo5xI
pDoiVhsLwg==
-----END CERTIFICATE-----

-----BEGIN CERTIFICATE-----
MIIF0zCCA7ugAwIBAgIIH8+hjWpIDREwDQYJKoZIhvcNAQELBQAwejELMAkGA1UE
BhMCQ04xFjAUBgNVBAoMDUFudCBGaW5hbmNpYWwxIDAeBgNVBAsMF0NlcnRpZmlj
YXRpb24gQXV0aG9yaXR5MTEwLwYDVQQDDChBbnQgRmluYW5jaWFsIENlcnRpZmlj
YXRpb24gQXV0aG9yaXR5IFIxMB4XDTE4MDMyMTEzNDg0MFoXDTM4MDIyODEzNDg0
MFowejELMAkGA1UEBhMCQ04xFjAUBgNVBAoMDUFudCBGaW5hbmNpYWwxIDAeBgNV
BAsMF0NlcnRpZmljYXRpb24gQXV0aG9yaXR5MTEwLwYDVQQDDChBbnQgRmluYW5j
aWFsIENlcnRpZmljYXRpb24gQXV0aG9yaXR5IFIxMIICIjANBgkqhkiG9w0BAQEF
AAOCAg8AMIICCgKCAgEAtytTRcBNuur5h8xuxnlKJetT65cHGemGi8oD+beHFPTk
rUTlFt9Xn7fAVGo6QSsPb9uGLpUFGEdGmbsQ2q9cV4P89qkH04VzIPwT7AywJdt2
xAvMs+MgHFJzOYfL1QkdOOVO7NwKxH8IvlQgFabWomWk2Ei9WfUyxFjVO1LVh0Bp
dRBeWLMkdudx0tl3+21t1apnReFNQ5nfX29xeSxIhesaMHDZFViO/DXDNW2BcTs6
vSWKyJ4YIIIzStumD8K1xMsoaZBMDxg4itjWFaKRgNuPiIn4kjDY3kC66Sl/6yTl
YUz8AybbEsICZzssdZh7jcNb1VRfk79lgAprm/Ktl+mgrU1gaMGP1OE25JCbqli1
Pbw/BpPynyP9+XulE+2mxFwTYhKAwpDIDKuYsFUXuo8t261pCovI1CXFzAQM2w7H
DtA2nOXSW6q0jGDJ5+WauH+K8ZSvA6x4sFo4u0KNCx0ROTBpLif6GTngqo3sj+98
SZiMNLFMQoQkjkdN5Q5g9N6CFZPVZ6QpO0JcIc7S1le/g9z5iBKnifrKxy0TQjtG
PsDwc8ubPnRm/F82RReCoyNyx63indpgFfhN7+KxUIQ9cOwwTvemmor0A+ZQamRe
9LMuiEfEaWUDK+6O0Gl8lO571uI5onYdN1VIgOmwFbe+D8TcuzVjIZ/zvHrAGUcC
AwEAAaNdMFswCwYDVR0PBAQDAgEGMAwGA1UdEwQFMAMBAf8wHQYDVR0OBBYEFF90
tATATwda6uWx2yKjh0GynOEBMB8GA1UdIwQYMBaAFF90tATATwda6uWx2yKjh0Gy
nOEBMA0GCSqGSIb3DQEBCwUAA4ICAQCVYaOtqOLIpsrEikE5lb+UARNSFJg6tpkf
tJ2U8QF/DejemEHx5IClQu6ajxjtu0Aie4/3UnIXop8nH/Q57l+Wyt9T7N2WPiNq
JSlYKYbJpPF8LXbuKYG3BTFTdOVFIeRe2NUyYh/xs6bXGr4WKTXb3qBmzR02FSy3
IODQw5Q6zpXj8prYqFHYsOvGCEc1CwJaSaYwRhTkFedJUxiyhyB5GQwoFfExCVHW
05ZFCAVYFldCJvUzfzrWubN6wX0DD2dwultgmldOn/W/n8at52mpPNvIdbZb2F41
T0YZeoWnCJrYXjq/32oc1cmifIHqySnyMnavi75DxPCdZsCOpSAT4j4lAQRGsfgI
kkLPGQieMfNNkMCKh7qjwdXAVtdqhf0RVtFILH3OyEodlk1HYXqX5iE5wlaKzDop
PKwf2Q3BErq1xChYGGVS+dEvyXc/2nIBlt7uLWKp4XFjqekKbaGaLJdjYP5b2s7N
1dM0MXQ/f8XoXKBkJNzEiM3hfsU6DOREgMc1DIsFKxfuMwX3EkVQM1If8ghb6x5Y
jXayv+NLbidOSzk4vl5QwngO/JYFMkoc6i9LNwEaEtR9PhnrdubxmrtM+RjfBm02
77q3dSWFESFQ4QxYWew4pHE0DpWbWy/iMIKQ6UZ5RLvB8GEcgt8ON7BBJeMc+Dyi
kT9qhqn+lw==
-----END CERTIFICATE-----

-----BEGIN CERTIFICATE-----
MIICiDCCAgygAwIBAgIIQX76UsB/30owDAYIKoZIzj0EAwMFADB6MQswCQYDVQQG
EwJDTjEWMBQGA1UECgwNQW50IEZpbmFuY2lhbDEgMB4GA1UECwwXQ2VydGlmaWNh
dGlvbiBBdXRob3JpdHkxMTAvBgNVBAMMKEFudCBGaW5hbmNpYWwgQ2VydGlmaWNh
dGlvbiBBdXRob3JpdHkgRTEwHhcNMTkwNDI4MTYyMDQ0WhcNNDkwNDIwMTYyMDQ0
WjB6MQswCQYDVQQGEwJDTjEWMBQGA1UECgwNQW50IEZpbmFuY2lhbDEgMB4GA1UE
CwwXQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkxMTAvBgNVBAMMKEFudCBGaW5hbmNp
YWwgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkgRTEwdjAQBgcqhkjOPQIBBgUrgQQA
IgNiAASCCRa94QI0vR5Up9Yr9HEupz6hSoyjySYqo7v837KnmjveUIUNiuC9pWAU
WP3jwLX3HkzeiNdeg22a0IZPoSUCpasufiLAnfXh6NInLiWBrjLJXDSGaY7vaokt
rpZvAdmjXTBbMAsGA1UdDwQEAwIBBjAMBgNVHRMEBTADAQH/MB0GA1UdDgQWBBRZ
4ZTgDpksHL2qcpkFkxD2zVd16TAfBgNVHSMEGDAWgBRZ4ZTgDpksHL2qcpkFkxD2
zVd16TAMBggqhkjOPQQDAwUAA2gAMGUCMQD4IoqT2hTUn0jt7oXLdMJ8q4vLp6sg
wHfPiOr9gxreb+e6Oidwd2LDnC4OUqCWiF8CMAzwKs4SnDJYcMLf2vpkbuVE4dTH
Rglz+HGcTLWsFs4KxLsq7MuU+vJTBUeDJeDjdA==
-----END CERTIFICATE-----

-----BEGIN CERTIFICATE-----
MIIDxTCCAq2gAwIBAgIUEMdk6dVgOEIS2cCP0Q43P90Ps5YwDQYJKoZIhvcNAQEF
BQAwajELMAkGA1UEBhMCQ04xEzARBgNVBAoMCmlUcnVzQ2hpbmExHDAaBgNVBAsM
E0NoaW5hIFRydXN0IE5ldHdvcmsxKDAmBgNVBAMMH2lUcnVzQ2hpbmEgQ2xhc3Mg
MiBSb290IENBIC0gRzMwHhcNMTMwNDE4MDkzNjU2WhcNMzMwNDE4MDkzNjU2WjBq
MQswCQYDVQQGEwJDTjETMBEGA1UECgwKaVRydXNDaGluYTEcMBoGA1UECwwTQ2hp
bmEgVHJ1c3QgTmV0d29yazEoMCYGA1UEAwwfaVRydXNDaGluYSBDbGFzcyAyIFJv
b3QgQ0EgLSBHMzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOPPShpV
nJbMqqCw6Bz1kehnoPst9pkr0V9idOwU2oyS47/HjJXk9Rd5a9xfwkPO88trUpz5
4GmmwspDXjVFu9L0eFaRuH3KMha1Ak01citbF7cQLJlS7XI+tpkTGHEY5pt3EsQg
wykfZl/A1jrnSkspMS997r2Gim54cwz+mTMgDRhZsKK/lbOeBPpWtcFizjXYCqhw
WktvQfZBYi6o4sHCshnOswi4yV1p+LuFcQ2ciYdWvULh1eZhLxHbGXyznYHi0dGN
z+I9H8aXxqAQfHVhbdHNzi77hCxFjOy+hHrGsyzjrd2swVQ2iUWP8BfEQqGLqM1g
KgWKYfcTGdbPB1MCAwEAAaNjMGEwHQYDVR0OBBYEFG/oAMxTVe7y0+408CTAK8hA
uTyRMB8GA1UdIwQYMBaAFG/oAMxTVe7y0+408CTAK8hAuTyRMA8GA1UdEwEB/wQF
MAMBAf8wDgYDVR0PAQH/BAQDAgEGMA0GCSqGSIb3DQEBBQUAA4IBAQBLnUTfW7hp
emMbuUGCk7RBswzOT83bDM6824EkUnf+X0iKS95SUNGeeSWK2o/3ALJo5hi7GZr3
U8eLaWAcYizfO99UXMRBPw5PRR+gXGEronGUugLpxsjuynoLQu8GQAeysSXKbN1I
UugDo9u8igJORYA+5ms0s5sCUySqbQ2R5z/GoceyI9LdxIVa1RjVX8pYOj8JFwtn
DJN3ftSFvNMYwRuILKuqUYSHc2GPYiHVflDh5nDymCMOQFcFG3WsEuB+EYQPFgIU
1DHmdZcz7Llx8UOZXX2JupWCYzK1XhJb+r4hK5ncf/w8qGtYlmyJpxk3hr1TfUJX
Yf4Zr0fJsGuv
-----END CERTIFICATE-----`

var alipayPublicCert = `-----BEGIN CERTIFICATE-----
MIIDsjCCApqgAwIBAgIQICUDEoNGkshFrVAKW0INLDANBgkqhkiG9w0BAQsFADCBgjELMAkGA1UE
BhMCQ04xFjAUBgNVBAoMDUFudCBGaW5hbmNpYWwxIDAeBgNVBAsMF0NlcnRpZmljYXRpb24gQXV0
aG9yaXR5MTkwNwYDVQQDDDBBbnQgRmluYW5jaWFsIENlcnRpZmljYXRpb24gQXV0aG9yaXR5IENs
YXNzIDIgUjEwHhcNMjUwMzEyMDY1MDEyWhcNMzAwMzExMDY1MDEyWjCBkjELMAkGA1UEBhMCQ04x
LTArBgNVBAoMJOWMl+S6rOaZuuWtkOaXtuepuuenkeaKgOaciemZkOWFrOWPuDEPMA0GA1UECwwG
QWxpcGF5MUMwQQYDVQQDDDrmlK/ku5jlrp0o5Lit5Zu9Kee9kee7nOaKgOacr+aciemZkOWFrOWP
uC0yMDg4OTMxOTA0NzI2NjEwMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxFo22xCl
D+gAJymP3s4jkIFPk9X9kQevg1hz0vPF4qqFwYLWoMnhQqfGLp4fcVEjVYfLfIxvAaDdZVEzMewb
qi1mGn9rJXjkmQ5YvANBnb6Egt7Hlmu3ECYelgLm7GBwotUrOxClhc2ejRvyA99ITbGY5dyOal6u
TI3g973wwrYfZVl3ARFEh5tRzEkpe2omR+zqY2pyl/QT4MYpMFXKHWNaH2gaxHH1sjajCbYTBjMn
JDU/MimGZ5ouVS1y3UOpgu4lHgFsQAnweoCCVL/44YNe7q+9DMCGiudOZ7A7ZQO/q2bNfciZdPUN
cPBGldSk8tftQ3n0tRO7Tn3Q+NHAXwIDAQABoxIwEDAOBgNVHQ8BAf8EBAMCA/gwDQYJKoZIhvcN
AQELBQADggEBAFOr6trKp6CP+LWSnicCWKIXWqBzlgeef3fYxG9yciaRsvEDyl5n7cofy198qw9I
vctV+NqIOPzmo29xnFceKZGGnlsxo8PTwtPnF1oFmPieQf0Q2czYio0e7T0lMckdjDwhIMBS5HXY
PpWfZoUnA5dku9Try+6D4pEDeY4/CJSPNkRI+m+uvn5Fva8NUuQNKRKIofcb4fpKPNu5E9vyfeyW
NTo3S1Hz/8nwKOGXOh3WQmbAOlOoeeoZt96AfcSQF1EH4tSDfwEBxuEY4mJnRRYHz9qCN4hi81QO
NIh8IsfT6GLy87crxYDcb10N5NJ+vx9g8JsBzgytJRHlmOU1sow=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIE4jCCAsqgAwIBAgIIYsSr5bKAMl8wDQYJKoZIhvcNAQELBQAwejELMAkGA1UEBhMCQ04xFjAU
BgNVBAoMDUFudCBGaW5hbmNpYWwxIDAeBgNVBAsMF0NlcnRpZmljYXRpb24gQXV0aG9yaXR5MTEw
LwYDVQQDDChBbnQgRmluYW5jaWFsIENlcnRpZmljYXRpb24gQXV0aG9yaXR5IFIxMB4XDTE4MDMy
MjE0MzQxNVoXDTM3MTEyNjE0MzQxNVowgYIxCzAJBgNVBAYTAkNOMRYwFAYDVQQKDA1BbnQgRmlu
YW5jaWFsMSAwHgYDVQQLDBdDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTE5MDcGA1UEAwwwQW50IEZp
bmFuY2lhbCBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eSBDbGFzcyAyIFIxMIIBIjANBgkqhkiG9w0B
AQEFAAOCAQ8AMIIBCgKCAQEAsLMfYaoRoPRbmDcAfXPCmKf43pWRN5yTXa/KJWO0l+mrgQvs89bA
NEvbDUxlkGwycwtwi5DgBuBgVhLliXu+R9CYgr2dXs8D8Hx/gsggDcyGPLmVrDOnL+dyeauheARZ
fA3du60fwEwwbGcVIpIxPa/4n3IS/ElxQa6DNgqxh8J9Xwh7qMGl0JK9+bALuxf7B541Gr4p0WEN
G8fhgjBV4w4ut9eQLOoa1eddOUSZcy46Z7allwowwgt7b5VFfx/P1iKJ3LzBMgkCK7GZ2kiLrL7R
iqV+h482J7hkJD+ardoc6LnrHO/hIZymDxok+VH9fVeUdQa29IZKrIDVj65THQIDAQABo2MwYTAf
BgNVHSMEGDAWgBRfdLQEwE8HWurlsdsio4dBspzhATAdBgNVHQ4EFgQUSqHkYINtUSAtDPnS8Xoy
oP9p7qEwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwDQYJKoZIhvcNAQELBQADggIB
AIQ8TzFy4bVIVb8+WhHKCkKNPcJe2EZuIcqvRoi727lZTJOfYy/JzLtckyZYfEI8J0lasZ29wkTt
a1IjSo+a6XdhudU4ONVBrL70U8Kzntplw/6TBNbLFpp7taRALjUgbCOk4EoBMbeCL0GiYYsTS0mw
7xdySzmGQku4GTyqutIGPQwKxSj9iSFw1FCZqr4VP4tyXzMUgc52SzagA6i7AyLedd3tbS6lnR5B
L+W9Kx9hwT8L7WANAxQzv/jGldeuSLN8bsTxlOYlsdjmIGu/C9OWblPYGpjQQIRyvs4Cc/mNhrh+
14EQgwuemIIFDLOgcD+iISoN8CqegelNcJndFw1PDN6LkVoiHz9p7jzsge8RKay/QW6C03KNDpWZ
EUCgCUdfHfo8xKeR+LL1cfn24HKJmZt8L/aeRZwZ1jwePXFRVtiXELvgJuM/tJDIFj2KD337iV64
fWcKQ/ydDVGqfDZAdcU4hQdsrPWENwPTQPfVPq2NNLMyIH9+WKx9Ed6/WzeZmIy5ZWpX1TtTolo6
OJXQFeItMAjHxW/ZSZTok5IS3FuRhExturaInnzjYpx50a6kS34c5+c8hYq7sAtZ/CNLZmBnBCFD
aMQqT8xFZJ5uolUaSeXxg7JFY1QsYp5RKvj4SjFwCGKJ2+hPPe9UyyltxOidNtxjaknOCeBHytOr
-----END CERTIFICATE-----`
var alipayAppId = "2021005126690253"

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
