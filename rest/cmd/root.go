package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zhlii/wechat-box/rest/internal/boot"
	"github.com/zhlii/wechat-box/rest/internal/callback"
	"github.com/zhlii/wechat-box/rest/internal/config"
	"github.com/zhlii/wechat-box/rest/internal/httpd"
	"github.com/zhlii/wechat-box/rest/internal/logs"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "rest",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg = "config.yaml"
		if cfgFile != "" {
			cfg = cfgFile
		}
		config.Init(cfg)

		logs.CreateLogger()

		boot := &boot.Boot{WcfPort: config.Data.Rpc.Port}
		err := boot.InitSDK()
		if err != nil {
			logs.Fatal(fmt.Sprintf("boot wx failed. err:%v", err))
		}

		defer boot.DestorySDK()

		if config.Data.Rpc.Enable {
			client := rpc.NewClient(config.Data.Rpc.Host, config.Data.Rpc.Port)

			if config.Data.Httpd.Enable {
				httpd := httpd.NewHttpServer(client)
				httpd.Start()
				defer httpd.Close()
			}

			err = client.Connect()
			if err != nil {
				logs.Fatal(fmt.Sprintf("connect wx rpc client failed. err:%v", err))
			}
			defer client.Close()

			for _, h := range callback.Setup() {
				handler := h
				err = client.RegisterCallback(func(msg *rpc.WxMsg) {
					handler.Callback(client, msg)
				})

				if err != nil {
					logs.Warn(fmt.Sprintf("register callback error:%v", err))
				}
			}
		}

		// 等待服务器停止信号
		chSig := make(chan os.Signal)
		signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
		<-chSig

		logs.Info("received quit signal.")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
