package main

import (
	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/registry"
	"github.com/liangdas/mqant/registry/consul"
	"github.com/nats-io/nats.go"
	"mqant-example/gate"
	"mqant-example/helloworld"
	"mqant-example/httpgateway"
	"mqant-example/rpctest"
	"mqant-example/web"
	"net/http"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	gate.JudgeGuest = func(session gate.Session) bool {
		if session.GetUserId() != "" {
			return false
		}
		return true
	}
	rs := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})

	nc, err := nats.Connect("nats://127.0.0.1:4222", nats.MaxReconnects(10000))
	if err != nil {
		log.Error("nats error %v", err)
		return
	}
	app := mqant.CreateApp(
		module.KillWaitTTL(1*time.Minute),
		module.Debug(true),  //只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		module.Nats(nc),     //指定nats rpc
		module.Registry(rs), //指定服务发现
		module.RegisterTTL(20*time.Second),
		module.RegisterInterval(10*time.Second),
	)
	err = app.Run( //模块都需要加到入口列表中传入框架
		helloworld.Module(),
		web.Module(),
		rpctest.Module(),
		mgate.Module(),
		httpgateway.Module(),
	)
	if err != nil {
		log.Error(err.Error())
	}
}
