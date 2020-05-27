/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package tabletest

import (
	"fmt"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/server"
	"time"
)

var Module = func() module.Module {
	this := new(tabletest)
	return this
}

type tabletest struct {
	basemodule.BaseModule
	room *room.Room
}

func (self *tabletest) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "tabletest"
}
func (self *tabletest) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *tabletest) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings,
		server.RegisterInterval(15*time.Second),
		server.RegisterTTL(30*time.Second),
	)
	self.room = room.NewRoom(self)
	self.GetServer().RegisterGO("HD_room_say", self.gatesay)
}

func (self *tabletest) Run(closeSig chan bool) {
	log.Info("%v模块运行中...", self.GetType())
	<-closeSig
}

func (self *tabletest) OnDestroy() {
	//一定别忘了继承
	self.BaseModule.OnDestroy()
}

func (self *tabletest) NewTable(module module.RPCModule, tableId string) (room.BaseTable, error) {
	table := NewTable(
		module,
		room.TableId(tableId),
		room.Router(func(TableId string) string {
			return fmt.Sprintf("%v://%v/%v", self.GetType(), self.GetServerId(), tableId)
		}),
		room.DestroyCallbacks(func(table room.BaseTable) error {
			log.Info("回收了房间: %v", table.TableId())
			_ = self.room.DestroyTable(table.TableId())
			return nil
		}),
	)
	return table, nil
}

func (self *tabletest) gatesay(session gate.Session, msg map[string]interface{}) (r string, err error) {
	table_id := msg["table_id"].(string)
	action := msg["action"].(string)
	table := self.room.GetTable(table_id)
	if table == nil {
		table, err = self.room.CreateById(self, table_id, self.NewTable)
		if err != nil {
			return "", err
		}
	}
	erro := table.PutQueue(action, session, msg)
	if erro != nil {
		return "", erro
	}
	return "success", nil
}
