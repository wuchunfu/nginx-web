package templatex

import (
	"github.com/valyala/fasttemplate"
	"github.com/wuchunfu/nginx-web/model/websiteModel"
	"github.com/wuchunfu/nginx-web/utils/datetimex"
	"net/url"
	"testing"
)

func TestReg(t *testing.T) {
	t.Run("template test", func(t *testing.T) {
		template := "https://{{host}}/?q={{query}}&foo={{bar}}{{bar}}"
		tempMap := map[string]interface{}{
			"host":  "google.com",
			"query": url.QueryEscape("hello=world"),
			"bar":   "foobar",
		}
		content := fasttemplate.New(template, "{{", "}}").ExecuteString(tempMap)
		t.Log(content)
	})

	t.Run("website template", func(t *testing.T) {
		website := new(websiteModel.Website)

		website.FileName = "test.conf"
		website.ServerName = "test.com"
		website.RootDirectory = "/tmp/www"
		website.HomePage = "index.html"
		website.HttpPort = 80
		website.SupportSsl = 0
		if website.SupportSsl == 1 {
			website.SupportSsl = 1
			website.HttpsPort = 443
			website.SslCertificate = "/tmp/tmp.cer"
			website.SslCertificateKey = "/tmp/tmp.key"
		}
		website.Status = 0
		website.CreateTime = datetimex.FormatNowDateTime()
		website.UpdateTime = datetimex.FormatNowDateTime()

		content := TemplateReplace(website)
		t.Log(string(content))
	})
}
