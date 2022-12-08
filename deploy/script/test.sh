#!/bin/sh

var="admin/api"
#imgName=${imgName/\//\-}
imgName=$(echo "${var}" | sed 's/\//\-/g')
echo "${imgName}"
srvName=${var%/*}
echo "$srvName"