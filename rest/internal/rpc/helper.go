package rpc

import (
	"compress/gzip"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// 解析数据库字段
// param field *DbField 字段
// return any 解析结果
func ParseDbField(field *DbField) any {
	str := string(field.Content)
	switch field.Type {
	case 1:
		n, _ := strconv.ParseInt(str, 10, 64)
		return n
	case 2:
		n, _ := strconv.ParseFloat(str, 64)
		return n
	case 4:
		return field.Content
	case 5:
		return nil
	default:
		return str
	}
}

func DownloadFile(str string) string {
	u, err := url.Parse(str)
	if err == nil && u.Scheme == "http" || u.Scheme == "https" {
		target := path.Join(os.TempDir(), strings.Trim(path.Base(u.Path), "/"))
		tmp, err := Download(str, target, false)
		if err == nil {
			time.AfterFunc(15*time.Minute, func() {
				os.RemoveAll(tmp)
			})
			return tmp
		}
	}
	return ""
}

func Download(url, target string, isGzip bool) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 默认读取器
	reader := resp.Body

	defer reader.Close()

	// 自动解压 gz 文件
	if isGzip || resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return "", err
		}
	}

	// 返回文件的名称
	return SaveStream(reader, target)

}

func SaveStream(reader io.Reader, target string) (string, error) {

	var err error
	var writer *os.File

	// 创建目标文件
	if target != "" {
		writer, err = os.Create(target)
	} else {
		writer, err = os.CreateTemp("", "tmp-*")
	}
	if err != nil {
		return "", err
	}
	defer writer.Close()

	// 写入文件数据
	_, err = io.Copy(writer, reader)
	if err != nil {
		return "", err
	}

	return writer.Name(), nil

}

func ContactType(wxid string) string {
	notFriends := map[string]string{
		"mphelper":    "公众平台助手",
		"fmessage":    "朋友推荐消息",
		"medianote":   "语音记事本",
		"floatbottle": "漂流瓶",
		"filehelper":  "文件传输助手",
		"newsapp":     "新闻",
	}
	if notFriends[wxid] != "" {
		return notFriends[wxid]
	}
	if strings.HasPrefix(wxid, "gh_") {
		return "公众号"
	}
	if strings.HasSuffix(wxid, "@chatroom") {
		return "群聊"
	}
	if strings.HasSuffix(wxid, "@openim") {
		return "企业微信"
	}
	return "好友"
}

func Rand(length uint) string {

	rs := make([]string, length)

	for i := uint(0); i < length; i++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rune(rand.Intn(26)+65)))
		} else {
			rs = append(rs, string(rune(rand.Intn(26)+97)))
		}
	}

	return strings.Join(rs, "")

}

func ToInt(str string) int {

	v, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return v

}
