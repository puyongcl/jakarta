#!/bin/sh

buildAll() {
  docker login reg.domain.com
  cd ../../ || exit

  for dir in "listener/rpc" "listener/mq" "im/api" "im/rpc" "im/mq" "order/rpc" "order/mq" "payment/api" "payment/rpc" "payment/mq" "mobile/api" "usercenter/rpc" "mqueue/job" "mqueue/scheduler" "admin/api" "chat/rpc" "chat/mq" "statistic/rpc" "bbs/rpc"
  do
    pwd
    buildImage "$dir"
  done
  return
}

build() {
  dir=$1
  docker login reg.domain.com
  cd ../../ || exit
  pwd

  buildImage "$dir"
}

buildImage() {
  ldir=$1
  echo buid "${ldir}"
  ver="latest"
  pwd
  imgName=$(echo "${ldir}" | sed 's/\//\-/g')
  docker build -f app/"${ldir}"/Dockerfile -t "${imgName}" .
  docker tag "${imgName}" reg.domain.com/bogota/"${imgName}":${ver}
  docker push reg.domain.com/bogota/"${imgName}":${ver}
  docker rmi "${imgName}"
  docker rmi reg.domain.com/bogota/"${imgName}":${ver}

  return
}