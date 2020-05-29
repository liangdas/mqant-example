// Copyright 2014 hey Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package table_test

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/liangdas/armyant/task"
	"github.com/liangdas/armyant/work"
	"github.com/liangdas/mqant/utils"
	"math/rand"
	"time"
)

func NewWork(manager *Manager) *Work {
	this := new(Work)
	this.manager = manager
	rand.NewSource(time.Now().UnixNano())
	this.name = fmt.Sprintf("%v", rand.Intn(100))
	//opts := this.GetDefaultOptions("tcp://127.0.0.1:3563")
	opts := this.GetDefaultOptions("ws://127.0.0.1:3653")
	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		fmt.Println("ConnectionLost", err.Error())
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		fmt.Println("OnConnectHandler")
	})
	// load root ca
	// 需要一个证书，这里使用的这个网站提供的证书https://curl.haxx.se/docs/caextract.html
	err := this.Connect(opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	this.On("/room/say", func(client MQTT.Client, msg MQTT.Message) {
		//服务端主动下发玩家加入事件
		fmt.Println("me is", this.name, msg.Topic(), "=》", string(msg.Payload()))
	})
	this.On("/room/join", func(client MQTT.Client, msg MQTT.Message) {
		//服务端主动下发玩家加入事件
		fmt.Println("me is", this.name, msg.Topic(), "=》", string(msg.Payload()))
	})
	return this
}

/**
Work 代表一个协程内具体执行任务工作者
*/
type Work struct {
	work.MqttWork
	manager *Manager
	name    string
}

func (this *Work) UnmarshalResult(payload []byte) map[string]interface{} {
	rmsg := map[string]interface{}{}
	json.Unmarshal(payload, &rmsg)
	return rmsg["Result"].(map[string]interface{})
}

/**
每一次请求都会调用该函数,在该函数内实现具体请求操作

task:=task.Task{
		N:1000,	//一共请求次数，会被平均分配给每一个并发协程
		C:100,		//并发数
		//QPS:10,		//每一个并发平均每秒请求次数(限流) 不填代表不限流
}

N/C 可计算出每一个Work(协程) RunWorker将要调用的次数
*/
func (this *Work) RunWorker(t task.Task) {
	msg, err := this.Request("tabletest/HD_room_say", []byte(
		mqanttools.Sprintf(`{"table_id":"{table_id}","action":"/room/join","name":"{name}"}`,
			map[string]string{"table_id": "mqantchat", "name": this.name},
		)))
	if err != nil {
		fmt.Println(msg.Topic(), err.Error())
		return
	}

	msg, err = this.Request("tabletest/HD_room_say", []byte(
		mqanttools.Sprintf(`{"table_id":"{table_id}","action":"/room/say","name":"{name}"}`,
			map[string]string{"table_id": "mqantchat", "name": this.name},
		)))
	if err != nil {
		fmt.Println(msg.Topic(), err.Error())
		return
	}

	fmt.Println(msg.Topic(), string(msg.Payload()))
}
func (this *Work) Init(t task.Task) {

}
func (this *Work) Close(t task.Task) {
	this.GetClient().Disconnect(0)
}
