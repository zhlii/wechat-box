package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zhlii/wechat-box/rest/internal/boot"
	"github.com/zhlii/wechat-box/rest/internal/logs"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "rest",
	Run: func(cmd *cobra.Command, args []string) {
		// var cfg = "config.yaml"
		// if cfgFile != "" {
		// 	cfg = cfgFile
		// }
		// config.Init(cfg)

		boot := &boot.Boot{WcfPort: 8888}
		err := boot.InitSDK()
		if err != nil {
			logs.Fatal(fmt.Sprintf("boot wx failed. err:%v", err))
		}

		defer boot.DestorySDK()

		client, err := rpc.NewClient("10.1.3.10:8888")
		if err != nil {
			logs.Fatal(fmt.Sprintf("create new wx rpc client failed. err:%v", err))
		}
		err = client.Connect()
		if err != nil {
			logs.Fatal(fmt.Sprintf("connect wx rpc client failed. err:%v", err))
		}

		client.RegisterCallback(func(msg *rpc.WxMsg) {
			fmt.Println(msg)
		})
		defer client.Close()

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
