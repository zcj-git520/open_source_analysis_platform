package pkg

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSendEmailWithTemplate(t *testing.T) {
	email := NewEmailSMTP(WithSmtpHost("smtp.qq.com"), WithSmtpPort(587), WithSmtpUsername("osap.work@qq.com"),
		WithSmtpPassword("dkcgahsacpdcbdjg"), WithFrom("osap.work@qq.com"),
		WithTo([]string{"zhaochengji@shds.cn"}))
	// 生成验证码
	rand.Seed(time.Now().UnixNano())
	verificationCode := rand.Intn(999999) + 100000

	// 邮件内容
	subject := "OSAP 注册验证码"
	body := fmt.Sprintf(`<td style="font-size:14px;color:#333;padding:24px 40px 0 40px">
                尊敬的用户您好！
                <br>
                <br>
                您的注册验证码是：<b>%d</b>，请在<b>5分钟内</b>进行验证, 过期将失效!
                <br> 
                如果该验证码不为您本人申请，请无视。
            </td>
`, verificationCode)
	//header := make(map[string]string)
	//header["From"] = from
	//header["To"] = "recipientEmail@example.com"
	//header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	//message := ""
	//for k, v := range header {
	//	message += fmt.Sprintf("%s: %s\r\n", k, v)
	//}
	//message += "\r\n" + body
	email.SendEmailSMTP(subject, body)
}
