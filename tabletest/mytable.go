package tabletest

import (
	"errors"
	"fmt"
	"github.com/liangdas/mqant-modules/room"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"reflect"
	"time"
)

type MyTable struct {
	room.QTable
	module  module.RPCModule
	players map[string]room.BasePlayer
}

func (this *MyTable) GetSeats() map[string]room.BasePlayer {
	return this.players
}
func (this *MyTable) GetModule() module.RPCModule {
	return this.module
}

func (this *MyTable) OnCreate() {
	//可以加载数据
	log.Info("MyTable OnCreate")
	//一定要调用QTable.OnCreate()
	this.QTable.OnCreate()
}

/**
每帧都会调用
*/
func (this *MyTable) Update(ds time.Duration) {

}

func NewTable(module module.RPCModule, opts ...room.Option) *MyTable {
	this := &MyTable{
		module:  module,
		players: map[string]room.BasePlayer{},
	}
	opts = append(opts, room.TimeOut(60))
	opts = append(opts, room.Update(this.Update))
	opts = append(opts, room.NoFound(func(msg *room.QueueMsg) (value reflect.Value, e error) {
		//return reflect.ValueOf(this.doSay), nil
		return reflect.Zero(reflect.ValueOf("").Type()), errors.New("no found handler")
	}))
	opts = append(opts, room.SetRecoverHandle(func(msg *room.QueueMsg, err error) {
		log.Error("Recover %v Error: %v", msg.Func, err.Error())
	}))
	opts = append(opts, room.SetErrorHandle(func(msg *room.QueueMsg, err error) {
		log.Error("Error %v Error: %v", msg.Func, err.Error())
	}))
	this.OnInit(this, opts...)
	this.Register("/room/say", this.doSay)
	this.Register("/room/join", this.doJoin)
	return this
}

func (this *MyTable) doSay(session gate.Session, msg map[string]interface{}) (err error) {
	player := this.FindPlayer(session)
	if player == nil {
		return errors.New("no join")
	}
	player.OnRequest(session)
	_ = this.NotifyCallBackMsg("/room/say", []byte(fmt.Sprintf("say hi from %v", msg["name"])))
	return nil
}

func (this *MyTable) doJoin(session gate.Session, msg map[string]interface{}) (err error) {
	player := &room.BasePlayerImp{}
	player.Bind(session)
	player.OnRequest(session)
	this.players[session.GetSessionId()] = player
	_ = this.NotifyCallBackMsg("/room/join", []byte(fmt.Sprintf("welcome to %v", msg["name"])))
	return nil
}
