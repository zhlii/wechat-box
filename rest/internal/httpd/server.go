package httpd

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/zhlii/wechat-box/rest/internal/config"
	"github.com/zhlii/wechat-box/rest/internal/httpd/midware"
	"github.com/zhlii/wechat-box/rest/internal/logs"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

type HttpServer struct {
	srv        *http.Server
	rpc_client *rpc.Client
}

func NewHttpServer(c *rpc.Client) *HttpServer {
	return &HttpServer{rpc_client: c}
}

func (s *HttpServer) Start() error {
	gin.SetMode((gin.ReleaseMode))

	go func() {
		engine := gin.New()
		if !config.Data.Common.IsProd {
			engine.Use(Cors())
		}

		r := engine.Group("/api")
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		r.Use(midware.OutputHandle)
		Route(s.rpc_client, r)

		// server a directory called static
		// ui static files
		engine.Use(SpaMiddleware("/", "./ui")) // your build of React or other SPA

		s.srv = &http.Server{
			Addr:    config.Data.Httpd.Addr,
			Handler: engine,
		}
		logs.Info(fmt.Sprintf("server is ready and listening on %s", config.Data.Httpd.Addr))
		err := s.srv.ListenAndServe()
		if err != http.ErrServerClosed {
			logs.Error(fmt.Sprintf("start http server error %v", err))
		}
	}()

	return nil
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization,X-Token,*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}

func SpaMiddleware(urlPrefix, spaDirectory string) gin.HandlerFunc {
	directory := static.LocalFile(spaDirectory, true)
	fileserver := http.FileServer(directory)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}

	return func(c *gin.Context) {
		if directory.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			c.Request.URL.Path = "/"
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func (s *HttpServer) Close() error {
	logs.Debug("close http server")
	return s.srv.Close()
}
