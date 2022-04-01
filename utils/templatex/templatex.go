package templatex

import (
	"github.com/valyala/fasttemplate"
	"github.com/wuchunfu/nginx-web/middleware/logx"
	"github.com/wuchunfu/nginx-web/model/websiteModel"
	"github.com/wuchunfu/nginx-web/template"
	"strconv"
)

func TemplateReplace(website *websiteModel.Website) []byte {
	content := ""
	tempMap := map[string]interface{}{}
	tempMap["httpPort"] = strconv.Itoa(int(website.HttpPort))
	tempMap["serverName"] = website.ServerName
	tempMap["rootDirectory"] = website.RootDirectory
	tempMap["homePage"] = website.HomePage
	tempMap["HTTP01PORT"] = strconv.Itoa(9180)
	if website.SupportSsl == 1 {
		tempMap["httpsPort"] = strconv.Itoa(int(website.HttpsPort))
		tempMap["sslCertificate"] = website.SslCertificate
		tempMap["sslCertificateKey"] = website.SslCertificateKey
		httpsTemplate := getTemplate("https-conf")
		content = fasttemplate.New(httpsTemplate, "{{ ", " }}").ExecuteString(tempMap)
	} else {
		httpTemplate := getTemplate("http-conf")
		content = fasttemplate.New(httpTemplate, "{{ ", " }}").ExecuteString(tempMap)
	}
	return []byte(content)
}

func getTemplate(fileName string) string {
	contentByte, err := template.TemplateFS.ReadFile(fileName)
	if err != nil {
		logx.GetLogger().Sugar().Errorf("Get template file failed: %s", err.Error())
		return ""
	}
	return string(contentByte)
}
