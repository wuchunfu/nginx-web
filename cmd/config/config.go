package config

import (
	"github.com/spf13/cobra"
	"github.com/wuchunfu/nginx-web/middleware/configx"
	"github.com/wuchunfu/nginx-web/run"
)

var StartCmd = &cobra.Command{
	Use:          "config",
	SilenceUsage: true,
	Short:        "Get Application config info",
	Example:      "goShortUrl config -f conf/config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		run.Run()
	},
}

func init() {
	cobra.OnInitialize(configx.InitConfig)

	setting := configx.ServerSetting

	StartCmd.PersistentFlags().StringVarP(&configx.ConfigFile, "configFile", "f", "conf/config.yaml", "config file")
	StartCmd.PersistentFlags().StringVar(&setting.Database.DbType, "dbType", "mysql", "database type")
	StartCmd.PersistentFlags().StringVar(&setting.Database.Host, "host", "127.0.0.1", "database host")
	StartCmd.PersistentFlags().IntVarP(&setting.Database.Port, "port", "p", 3306, "database port")
	StartCmd.PersistentFlags().StringVar(&setting.Database.DbName, "dbName", "", "database name")
	StartCmd.PersistentFlags().StringVar(&setting.Database.Username, "username", "", "database username")
	StartCmd.PersistentFlags().StringVar(&setting.Database.Password, "password", "", "database password")

	// 必须配置项
	_ = StartCmd.MarkFlagRequired("configFile")

	// 使用viper可以绑定flag
	_ = configx.Vip.BindPFlag("database.dbType", StartCmd.PersistentFlags().Lookup("dbType"))
	_ = configx.Vip.BindPFlag("database.host", StartCmd.PersistentFlags().Lookup("host"))
	_ = configx.Vip.BindPFlag("database.port", StartCmd.PersistentFlags().Lookup("port"))
	_ = configx.Vip.BindPFlag("database.dbName", StartCmd.PersistentFlags().Lookup("dbName"))
	_ = configx.Vip.BindPFlag("database.username", StartCmd.PersistentFlags().Lookup("username"))
	_ = configx.Vip.BindPFlag("database.password", StartCmd.PersistentFlags().Lookup("password"))

	// 设置默认值
	configx.Vip.SetDefault("system.port", "9090")
}
