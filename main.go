package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/queue"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/threading"
	"tt90.cc/ucenter/internal/config"
	"tt90.cc/ucenter/internal/handler"
	"tt90.cc/ucenter/internal/jobs"
	"tt90.cc/ucenter/internal/middleware"
	"tt90.cc/ucenter/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/main.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 请求日志
	server.Use(middleware.NewRequestLogMiddleware(c.Auth.AccessSecret).Handle)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	if c.Mode == service.TestMode || c.Mode == service.ProMode {
		threading.GoSafe(func() {
			jobs.RegisterJobs(ctx)
		})

		threading.GoSafe(func() {
			queue.NewQueue(ctx, nil).Start()
		})
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
