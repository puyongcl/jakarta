#!/bin/sh

#api
goctl api go -api server.api -dir ../  --style=goZero

# swagger
goctl api plugin -plugin goctl-swagger="swagger -filename chat.json" -api chat.api -dir .
goctl api plugin -plugin goctl-swagger="swagger -filename listener.json" -api listener.api -dir .
goctl api plugin -plugin goctl-swagger="swagger -filename order.json" -api order.api -dir .
#goctl api plugin -plugin goctl-swagger="swagger -filename payment.json" -api payment.api -dir .
goctl api plugin -plugin goctl-swagger="swagger -filename user.json" -api user.api -dir .
goctl api plugin -plugin goctl-swagger="swagger -filename stat.json" -api stat.api -dir .
goctl api plugin -plugin goctl-swagger="swagger -filename bbs.json" -api bbs.api -dir .
