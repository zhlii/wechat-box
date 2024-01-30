package midware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func OutputHandle(c *gin.Context) {

	c.Next()

	// 输出错误信息

	if err, ok := c.Get("Error"); ok {
		c.AbortWithStatusJSON(exitCode(c, 400), newErrorMessage(err))
		return
	}

	// 输出请求结果

	msg := c.GetString("Message")

	if res, ok := c.Get("Data"); ok || msg != "" {
		data := newPayload(res, msg)
		c.AbortWithStatusJSON(exitCode(c, 200), data)
		return
	}

	// 输出HTML内容

	if htm := c.GetString("HTML"); htm != "" {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, htm)
		c.Abort()
		return
	}

	// 捕获异常返回

	c.AbortWithStatusJSON(500, newErrorMessage("内部错误"))

}

func exitCode(c *gin.Context, code int) int {

	if code := c.GetInt("ExitCode"); code > 100 {
		return code
	}

	return code

}

func newError(data any) error {

	if err, ok := data.(error); ok {
		return err
	}

	if err, ok := data.(string); ok {
		return errors.New(err)
	}

	return errors.New("未知错误")

}

func newErrorMessage(data any) gin.H {

	if err, ok := data.(error); ok {
		return gin.H{"Error": gin.H{"Message": err.Error()}}
	}

	if err, ok := data.(string); ok {
		return gin.H{"Error": gin.H{"Message": err}}
	}

	return gin.H{"Error": data}

}

func newPayload(data any, msg string) gin.H {

	payload := gin.H{"Data": data}

	if msg != "" {
		payload["Message"] = msg
	}

	return payload
}
