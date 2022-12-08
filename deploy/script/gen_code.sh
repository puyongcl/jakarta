# 生成api业务代码 进入"jakarta/app/xx/api/desc"目录下，执行下面命令
# goctl api go -api *.api -dir ../  --style=goZero
# 生成rpc业务代码
#【注】需要按照go-zero说明文档安装环境
# 进入"jakarta/app/xx/rpc/pb"目录下，执行下面命令
# goctl rpc protoc *.proto --go_out=../ --go-grpc_out=../  --zrpc_out=../ --style=goZero
# 去除proto中的json的omitempty
# mac: sed -i "" 's/,omitempty//g' *.pb.go
# linux: sed -i 's/,omitempty//g' *.pb.go

# 数据库生成pb文件
#mysql
# sql2pb -go_package ./pb -host 192.168.1.12 -package pb -password PXDN93VRKUm8TeE7 -port 33069 -schema jakarta -service_name listener -user root > listener.proto
#postgresql不可用
# sql2pb -db pg -go_package ./pb -host 192.168.1.12 -package pb -password postgres -port 5432 -schema jakarta_listener -service_name listener -user jakarta > listener.proto

# 生成md文档
#goctl api doc -dir ./usercenter.api -o ./

# 生成swagger文档
# https://github.com/zeromicro/goctl-swagger
# goctl api plugin -plugin goctl-swagger="swagger -filename usercenter.json" -api usercenter.api -dir .