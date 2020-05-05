/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package rpctest

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/rpc/pb"
	"github.com/liangdas/mqant/server"
	"time"
)

var Module = func() module.Module {
	this := new(rpctest)
	return this
}

type rpctest struct {
	basemodule.BaseModule
}

func (self *rpctest) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "rpctest"
}
func (self *rpctest) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *rpctest) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings,
		server.RegisterInterval(15*time.Second),
		server.RegisterTTL(30*time.Second),
	)
	self.GetServer().RegisterGO("/test/proto", self.testProto)
	self.GetServer().RegisterGO("/test/marshal", self.testMarshal)
}

func (self *rpctest) Run(closeSig chan bool) {
	log.Info("%v模块运行中...", self.GetType())
	<-closeSig
}

func (self *rpctest) OnDestroy() {
	//一定别忘了关闭RPC
	_ = self.GetServer().OnDestroy()
}
func (self *rpctest) testProto(req *rpcpb.ResultInfo) (*rpcpb.ResultInfo, error) {
	r := &rpcpb.ResultInfo{Error: *proto.String(fmt.Sprintf("你说: %v", req.Error))}
	return r, nil
}

func (self *rpctest) testMarshal(req Req) (*Rsp, error) {
	r := &Rsp{Msg: fmt.Sprintf("你的ID：%v", req.Id)}
	return r, nil
}
