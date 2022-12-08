#!/bin/sh

. ./func.sh
buildAll

# docker build -f app/admin/api/Dockerfile -t admin-api .
# docker tag admin-api reg.domain.com/bogota/admin-api:latest
# docker push reg.domain.com/bogota/admin-api:latest
# docker rmi admin-api
# docker rmi reg.domain.com/bogota/admin-api:latest