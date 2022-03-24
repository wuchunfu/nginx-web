package main

import (
	"flag"
	"github.com/wuchunfu/nginx-web/cmd"
	"github.com/wuchunfu/nginx-web/middleware/configx"
	"github.com/wuchunfu/nginx-web/run"
)

func main() {
	dev := flag.Bool("dev", false, "Is it a dev env")
	flag.Parse()
	args := flag.Args()
	if *dev {
		if len(args) > 0 {
			configx.ConfigFile = args[0]
		} else {
			configx.ConfigFile = "conf/config.yaml"
		}
		configx.InitConfig()
		run.Run()
	} else {
		cmd.Execute()
	}
}
