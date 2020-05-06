/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package web

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/registry"
	"github.com/liangdas/mqant/rpc"
	"github.com/liangdas/mqant/rpc/pb"
	"github.com/liangdas/mqant/selector"
	"io"
	"math/rand"
	"mqant-example/rpctest"
	"net/http"
	"sync"
	"time"
)

var Module = func() module.Module {
	this := new(Web)
	return this
}

type Web struct {
	basemodule.BaseModule
}

func (self *Web) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Web"
}
func (self *Web) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *Web) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
}

func (self *Web) startHttpServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
		rstr, err := mqrpc.String(
			self.RpcCall(
				ctx,
				"helloworld", //要访问的moduleType
				"/say/hi",    //访问模块中handler路径
				mqrpc.Param(r.Form.Get("name")),
				selector.WithStrategy(func(services []*registry.Service) selector.Next {
					var nodes []*registry.Node

					// Filter the nodes for datacenter
					for _, service := range services {
						if service.Version != "1.0.0" {
							continue
						}
						for _, node := range service.Nodes {
							nodes = append(nodes, node)
							if node.Metadata["state"] == "alive" || node.Metadata["state"] == "" {
								nodes = append(nodes, node)
							}
						}
					}

					var mtx sync.Mutex
					//log.Info("services[0] $v",services[0].Nodes[0])
					return func() (*registry.Node, error) {
						mtx.Lock()
						defer mtx.Unlock()
						if len(nodes) == 0 {
							return nil, fmt.Errorf("no node")
						}
						index := rand.Intn(int(len(nodes)))
						return nodes[index], nil
					}
				}),
			),
		)

		//SvrSession, err :=self.GetRouteServer("helloworld",
		//	selector.WithStrategy(func(services []*registry.Service) selector.Next {
		//		var nodes []*registry.Node
		//
		//		// Filter the nodes for datacenter
		//		for _, service := range services {
		//			for _, node := range service.Nodes {
		//              nodes = append(nodes, node)
		//				if node.Metadata["version"] == "1.0.0" {
		//					nodes = append(nodes, node)
		//				}
		//			}
		//		}
		//
		//		var mtx sync.Mutex
		//		//log.Info("services[0] $v",services[0].Nodes[0])
		//		return func() (*registry.Node, error) {
		//			mtx.Lock()
		//			defer mtx.Unlock()
		//			if len(nodes) == 0 {
		//				return nil, fmt.Errorf("no node")
		//			}
		//			index := rand.Intn(int(len(nodes)))
		//			return nodes[index], nil
		//		}
		//	}), )
		//if err != nil {
		//	log.Warning("HelloWorld error:%v", err.Error())
		//	return
		//}
		//rstr, err :=mqrpc.String(SvrSession.Call(ctx, "/say/hi", r.Form.Get("name")))
		//if err != nil {
		//	log.Warning("HelloWorld error:%v", err)
		//	return
		//}

		log.Info("RpcCall %v , err %v", rstr, err)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
		}
		_, _ = io.WriteString(w, rstr)
	})

	http.HandleFunc("/test/proto", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
		protobean := new(rpcpb.ResultInfo)
		err := mqrpc.Proto(protobean, func() (reply interface{}, errstr interface{}) {
			return self.RpcCall(
				ctx,
				"rpctest",     //要访问的moduleType
				"/test/proto", //访问模块中handler路径
				mqrpc.Param(&rpcpb.ResultInfo{Error: *proto.String(r.Form.Get("message"))}),
			)
		})
		log.Info("RpcCall %v , err %v", protobean, err)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
		}
		_, _ = io.WriteString(w, protobean.Error)
	})

	http.HandleFunc("/test/marshal", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
		rspbean := new(rpctest.Rsp)
		err := mqrpc.Marshal(rspbean, func() (reply interface{}, errstr interface{}) {
			return self.RpcCall(
				ctx,
				"rpctest",       //要访问的moduleType
				"/test/marshal", //访问模块中handler路径
				mqrpc.Param(&rpctest.Req{Id: r.Form.Get("mid")}),
			)
		})
		log.Info("RpcCall %v , err %v", rspbean, err)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
		}
		_, _ = io.WriteString(w, rspbean.Msg)
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Info("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	// returning reference so caller can call Shutdown()
	return srv
}

func (self *Web) Run(closeSig chan bool) {
	log.Info("web: starting HTTP server :8080")
	srv := self.startHttpServer()
	<-closeSig
	log.Info("web: stopping HTTP server")
	// now close the server gracefully ("shutdown")
	// timeout could be given instead of nil as a https://golang.org/pkg/context/
	if err := srv.Shutdown(nil); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	log.Info("web: done. exiting")
}

func (self *Web) OnDestroy() {
	//一定别忘了关闭RPC
	_ = self.GetServer().OnDestroy()
}
