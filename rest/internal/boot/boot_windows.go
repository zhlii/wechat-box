//go:build windows
// +build windows

package boot

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/zhlii/wechat-box/rest/internal/logs"

	"golang.org/x/sys/windows"
)

type Boot struct {
	SdkFilePath string // sdk.dll文件的路径
	WcfPort     int    // wcf 监听的端口
}

func (b *Boot) InitSDK() error {
	executable, _ := os.Executable()
	b.SdkFilePath = path.Join(filepath.Dir(executable), "sdk.dll")

	logs.Debug(fmt.Sprintf("sdk file path: %s", b.SdkFilePath))

	callFuncInDll(b.SdkFilePath, "WxInitSDK", uintptr(1), uintptr(b.WcfPort))

	time.Sleep(time.Second * 10)

	return nil
}

func (b *Boot) DestorySDK() error {
	if b.SdkFilePath == "" {
		return nil
	}

	err := callFuncInDll(b.SdkFilePath, "WxDestroySDK", uintptr(0))
	if err != nil {
		return err
	}

	return nil
}

func callFuncInDll(dll_file string, fn string, a ...uintptr) error {
	sdk, err := windows.LoadDLL(dll_file)
	if err != nil {
		logs.Error(fmt.Sprintf("failed to load %s", dll_file))
		return err
	}
	logs.Debug("dll loaded.")

	defer sdk.Release()

	// 查找 fn 函数
	proc, err := sdk.FindProc(fn)
	if err != nil {
		logs.Error(fmt.Sprintf("failed to call %s in dll %s", fn, dll_file))
		return err
	}

	logs.Debug("proc finded.")

	// 执行 fn(a...)
	r1, r2, err := proc.Call(a...)
	logs.Debug(fmt.Sprintf("call dll:%s r1:%s r2:%s", fn, r1, r2))
	return err
}
