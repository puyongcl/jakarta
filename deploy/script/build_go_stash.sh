#!/bin/sh

docker login reg.domain.com
ver="latest"
imgName="go-stash"
docker build -f stash/Dockerfile -t "${imgName}" .
docker tag "${imgName}" reg.domain.com/bogota/"${imgName}":${ver}
docker push reg.domain.com/bogota/"${imgName}":${ver}
docker rmi "${imgName}"
docker rmi reg.domain.com/bogota/"${imgName}":${ver}