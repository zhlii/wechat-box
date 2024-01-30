//go:build !windows
// +build !windows

package boot

type Boot struct {
	SdkFilePath string // sdk.dll文件的路径
	WcfPort     int    // wcf 监听的端口
}

func (b *Boot) InitSDK() error {
	return nil
}

func (b *Boot) DestorySDK() error {
	return nil
}
