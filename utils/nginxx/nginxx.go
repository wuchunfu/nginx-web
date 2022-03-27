package nginxx

import (
	"errors"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Test() error {
	out, err := exec.Command("nginx", "-t").CombinedOutput()
	if err != nil {
		return err
	}

	output := string(out)
	if strings.Contains(output, "failed") {
		return errors.New(output)
	}
	return nil
}

func Reload() (string, error) {
	out, err := exec.Command("nginx", "-s", "reload").CombinedOutput()
	if err != nil {
		return "", err
	}

	output := string(out)
	return output, nil
}

func GetConfPath(dir string) (string, error) {
	out, err := exec.Command("nginx", "-V").CombinedOutput()
	if err != nil {
		return "", err
	}

	reg, _ := regexp.Compile("--conf-path=(.*)/(.*.conf)")
	confPath := reg.FindStringSubmatch(string(out))[1]
	filePath := filepath.Join(confPath, dir)
	return filePath, nil
}
