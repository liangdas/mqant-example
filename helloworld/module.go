/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package helloworld

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
)

var Module = func() module.Module {
	this := new(HellWorld)
	return this
}

type HellWorld struct {
	basemodule.BaseModule
}

func (self *HellWorld) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "helloworld"
}
func (self *HellWorld) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *HellWorld) OnAppConfigurationLoaded(app module.App) {
	//当App初始化时调用，这个接口不管这个模块是否在这个进程运行都会调用
	self.BaseModule.OnAppConfigurationLoaded(app)
}
func (self *HellWorld) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.GetServer().Options().Metadata["state"] = "alive"
	self.GetServer().RegisterGO("/say/hi", self.say) //handler
	self.GetServer().RegisterGO("HD_say", self.gatesay)
	log.Info("%v模块初始化完成...", self.GetType())
}

func (self *HellWorld) Run(closeSig chan bool) {
	log.Info("%v模块运行中...", self.GetType())
	<-closeSig
	log.Info("%v模块已停止...", self.GetType())
}

func (self *HellWorld) OnDestroy() {
	//一定继承
	self.BaseModule.OnDestroy()
	log.Info("%v模块已回收...", self.GetType())
}
func (self *HellWorld) say(name string) (r string, err error) {
	return fmt.Sprintf("hi %v", name), nil
}

func (self *HellWorld) gatesay(session gate.Session, msg map[string]interface{}) (r string, err error) {
	session.Send("/gate/send/test", []byte(fmt.Sprintf("send hi to %v", msg["name"])))
	return fmt.Sprintf("hi %v 你在网关 %v", msg["name"], session.GetServerId()), nil
}
