package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TransResult struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Result    []*Result `json:"trans_result"`
	ErrorCode string    `json:"error_code"`
	ErrorMsg  string    `json:"error_msg"`
}

type Result struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

// 申请的信息
var AppID = "20230303001583911"
var SecretKey = "8VVLEVnByjAq0vaEScDX"

const (
	Url = "http://api.fanyi.baidu.com/api/trans/vip/translate"
)

type TranslateModel struct {
	From      string
	To        string
	Appid     string
	SecretKey string
	Salt      int
	Sign      string
}

func NewTranslateModeler(appid, secretKey, from, to string) *TranslateModel {
	return &TranslateModel{
		SecretKey: secretKey,
		Appid:     appid,
		From:      from,
		To:        to,
		Salt:      time.Now().Second(),
	}
}

func (t *TranslateModel) toValues(content string) url.Values {
	values := url.Values{
		"q":     {content},
		"from":  {t.From},
		"to":    {t.To},
		"appid": {t.Appid},
		"salt":  {strconv.Itoa(t.Salt)},
		"sign":  {t.sumString(content)},
	}
	return values
}

// 计算文本的md5值
func (t *TranslateModel) sumString(q string) string {
	content := t.Appid + q + strconv.Itoa(t.Salt) + t.SecretKey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(content))
	bys := md5Ctx.Sum(nil)
	value := hex.EncodeToString(bys)
	return value
}

func (t *TranslateModel) Translate(content string) string {
	values := t.toValues(content)
	resp, err := http.PostForm(Url, values)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var ts TransResult
	_ = json.Unmarshal(body, &ts)
	if ts.ErrorCode != "" {
		//return ts.ErrorMsg
		return ""
	}
	if len(ts.Result) > 0 {
		return ts.Result[0].Dst
	}
	return ""
}
