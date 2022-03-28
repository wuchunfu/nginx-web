package nginxx

import (
	"os/exec"
	"regexp"
	"testing"
)

func TestGetNginx(t *testing.T) {
	out, err := exec.Command("nginx", "-V").CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s\n", out)

	reg, _ := regexp.Compile("--conf-path=(.*)/(.*.conf)")
	str := reg.FindStringSubmatch(string(out))[1]
	t.Log(str)
}
