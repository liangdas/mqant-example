/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package helloworld

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
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
	return "HelloWorld"
}
func (self *HellWorld) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *HellWorld) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.GetServer().Options().Metadata["state"] = "alive"
	self.GetServer().RegisterGO("/say/hi", self.say) //handler
	log.Info("%v模块初始化完成...", self.GetType())
}

func (self *HellWorld) Run(closeSig chan bool) {
	log.Info("%v模块运行中...", self.GetType())
	<-closeSig
	log.Info("%v模块已停止...", self.GetType())
}

func (self *HellWorld) OnDestroy() {
	//一定别忘了关闭RPC
	_ = self.GetServer().OnDestroy()
	log.Info("%v模块已回收...", self.GetType())
}
func (self *HellWorld) say(name string) (r string, err error) {
	return fmt.Sprintf("hi %v", name), nil
}
