package callback

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zhlii/wechat-box/rest/internal/config"
	"github.com/zhlii/wechat-box/rest/internal/logs"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

func handlerSpark() {
	handlers["spark"] = &Handler{
		Callback: func(c *rpc.Client, msg *rpc.WxMsg) {
			cfg := config.Data.Callbacks["spark"]
			switch msg.Type {
			case 1:
				if msg.IsGroup {
					// answer, err := Invoke(msg.Roomid, msg.Content, cfg["appid"], cfg["apiKey"], cfg["apiSecret"])
					// if err != nil {
					// 	logs.Error(fmt.Sprintf("call spark error: %v", err))
					// } else {
					// 	c.CmdClient.SendTxt(answer, msg.Roomid, "")
					// }
				} else if msg.IsSelf && msg.Receiver == c.Usr.Wxid {
					answer, err := spark_ask(msg.Receiver, msg.Content, cfg["appid"], cfg["apiKey"], cfg["apiSecret"])
					if err != nil {
						logs.Error(fmt.Sprintf("call spark error: %v", err))
					} else {
						c.CmdClient.SendTxt(answer, msg.Receiver, "")
					}
				}
			}
		},
	}
}

var msgHistories = make(map[string][]Message)

func AppendHistory(id string, items ...Message) {
	if len(msgHistories[id]) >= 10 {
		msgHistories[id] = msgHistories[id][len(items):]
	}

	msgHistories[id] = append(msgHistories[id], items...)
}

/**
 *  WebAPI 接口调用示例 接口文档（必看）：https://www.xfyun.cn/doc/spark/Web.html
 * 错误码链接：https://www.xfyun.cn/doc/spark/%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E.html（code返回错误码时必看）
 * @author iflytek
 */

var (
	hostUrl = "wss://spark-api.xf-yun.com/v3.1/chat"
)

func spark_ask(sender, text, appid, apiKey, apiSecret string) (string, error) {
	if len(text) == 0 {
		return "", errors.New("text is empty.")
	}

	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(assembleAuthUrl1(hostUrl, apiKey, apiSecret), nil)
	if err != nil {
		return "", errors.New(readResp(resp) + err.Error())
	} else if resp.StatusCode != 101 {
		return "", errors.New(readResp(resp) + err.Error())
	}

	go func() {
		data := genParams1(appid, sender, text)
		conn.WriteJSON(data)
	}()

	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("read message error: %v", err)
		}

		var data map[string]interface{}
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			return "", fmt.Errorf("Error parsing JSON: %v", err)
		}

		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			return "", fmt.Errorf("code!=0, %v", data["payload"])
		}

		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			answer += content
		} else {
			answer += content
			// usage := payload["usage"].(map[string]interface{})
			// temp := usage["text"].(map[string]interface{})
			// totalTokens := temp["total_tokens"].(float64)
			conn.Close()
			break
		}
	}

	AppendHistory(sender, Message{Role: "assistant", Content: answer})

	return answer, nil
}

// 生成参数
func genParams1(appid, sender, question string) map[string]interface{} {

	message := Message{Role: "user", Content: question}
	AppendHistory(sender, message)

	for _, v := range msgHistories[sender] {
		fmt.Printf("history: %v\n", v)
	}

	data := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": appid,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      "generalv3",
				"temperature": float64(0.8),
				"max_tokens":  int64(8192),
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": msgHistories[sender],
			},
		},
	}
	return data
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
