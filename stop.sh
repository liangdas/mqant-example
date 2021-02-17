#!/bin/sh
sp_pid=`ps -ef | grep mqant-example | grep -v grep | awk '{print $2}'`
if [ -z "$sp_pid" ];
then
 echo "[ not find mqant-example pid ]"
else
 echo "find result: $sp_pid "
 kill -9 $sp_pid
fi