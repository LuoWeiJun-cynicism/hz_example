// Code generated by hertz generator.

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/arl/statsviz"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/hertz-contrib/cors"
)

func main() {
	// 自定义扩展
	h := server.Default(
		server.WithHostPorts("0.0.0.0:1111"),
		server.WithMaxRequestBodySize(20<<20),
		// 配置网络模型
		server.WithTransport(standard.NewTransporter),
	)
	h.Static("file", "./")
	h.StaticFS("/", &app.FS{})
	// like SimpleHTTPServer
	h.StaticFS("/try_dir", &app.FS{Root: "./", GenerateIndexPages: true, PathRewrite: app.NewPathSlashesStripper(1)})
	h.StaticFile("/main", "./file/staticFile/main.go")
	// 配置跨域中间件
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"}, // Allowed domains, need to bring schema
		AllowMethods:     []string{"PUT", "PATCH"},    // Allowed request methods
		AllowHeaders:     []string{"Origin"},          // Allowed request headers
		ExposeHeaders:    []string{"Content-Length"},  // Request headers allowed in the upload_file
		AllowCredentials: true,                        // Whether cookies are attached
		AllowOriginFunc: func(origin string) bool { // Custom domain detection with lower priority than AllowOrigins
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour, // Maximum length of upload_file-side cache preflash requests (seconds)
	}))

	//h.Use(basic_auth.BasicAuth(map[string]string{
	//	"test1": "value1",
	//	"test2": "value2",
	//}))
	// 使用自定义中间件
	//h.Use(middleware.MyMiddleware2())
	//
	// 在线指标监控
	statsviz.RegisterDefault()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	register(h)
	h.Spin()
}