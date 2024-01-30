package helper

import (
	"errors"
	"strings"

	"github.com/kbinani/screenshot"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
)

func WxLoginQrcode() (string, error) {
	// 获取第一个显示器的屏幕截图
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", err
	}
	// 将图片转换为 BinaryBitmap
	source := gozxing.NewLuminanceSourceFromImage(img)
	bmp, _ := gozxing.NewBinaryBitmap(gozxing.NewHybridBinarizer(source))
	// 检测图片中的多个二维码
	qrReader := qrcode.NewQRCodeMultiReader()
	results, err := qrReader.DecodeMultipleWithoutHint(bmp)
	if err != nil {
		return "", err
	}
	// 挑出微信登录的二维码
	for _, result := range results {
		url := result.String()
		if strings.HasPrefix(url, "http://weixin.qq.com/x/") {
			return url, nil
		}
	}
	return "", errors.New("未找到二维码")
}
