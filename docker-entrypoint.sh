#!/bin/sh
echo "starting nats-server"
/nats-server >/app/nats.log &
echo "starting consul"
/consul agent --dev > /app/consul.log &
echo "sleep 3s"
sleep 3s #暂停3s等待nats consul启动完成
echo "starting mqant-example"
sh /app/start.sh