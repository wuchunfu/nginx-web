package nginxx

import (
	"errors"
	"github.com/wuchunfu/nginx-web/middleware/logx"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Test() error {
	out, err := exec.Command("nginx", "-t").CombinedOutput()
	if err != nil {
		logx.GetLogger().Sugar().Error(err)
	}

	output := string(out)
	logx.GetLogger().Sugar().Info(output)

	if strings.Contains(output, "failed") {
		return errors.New(output)
	}
	return nil
}

func Reload() string {
	out, err := exec.Command("nginx", "-s", "reload").CombinedOutput()
	if err != nil {
		logx.GetLogger().Sugar().Error(err)
	}

	output := string(out)
	logx.GetLogger().Sugar().Info(output)
	return output
}

func GetConfPath(dir string) string {
	out, err := exec.Command("nginx", "-V").CombinedOutput()
	if err != nil {
		logx.GetLogger().Sugar().Error(err)
		return ""
	}

	reg, _ := regexp.Compile("--conf-path=(.*)/(.*.conf)")
	confPath := reg.FindStringSubmatch(string(out))[1]
	filePath := filepath.Join(confPath, dir)
	return filePath
}
