#!/bin/sh

BASE_DIR=$(pwd)
echo $BASE_DIR

for var in "listener/rpc" "listener/mq" "im/api" "im/rpc" "im/mq" "order/rpc" "order/mq" "payment/api" "payment/rpc" "payment/mq" "mobile/api" "usercenter/rpc" "mqueue/job" "mqueue/scheduler" "admin/api" "chat/rpc" "chat/mq" "statistic/rpc" "bbs/rpc"
do
  pwd
  BUILD_DIR="../../app/"$var
  cd $BUILD_DIR || exit
  pwd
  go build
  go clean

  cd $BASE_DIR || exit
done