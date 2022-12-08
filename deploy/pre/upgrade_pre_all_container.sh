#!/bin/sh


. ./func.sh

for dir in "listener/rpc" "listener/mq" "im/api" "im/rpc" "im/mq" "order/rpc" "order/mq" "payment/api" "payment/rpc" "payment/mq" "mobile/api" "usercenter/rpc" "mqueue/job" "mqueue/scheduler" "admin/api" "chat/rpc" "chat/mq" "statistic/rpc" "bbs/rpc"
do
  upgrade_pre "${dir}"
done