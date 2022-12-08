# 在model目录下执行 生成CRUD+cache代码
#mysql
# goctl model mysql ddl -c -src usercenter.sql -dir . --style=goZero

#postgresql
#goctl model pg datasource -url="postgres://jakarta:postgres@192.168.1.12:5432/jakarta_listener?sslmode=disable" -schema="jakarta" -table="*" -c -dir="." --style=goZero

# payment flow不要用cache
#goctl model pg datasource -url="postgres://jakarta:postgres@192.168.1.12:5432/jakarta_payment?sslmode=disable" -schema="jakarta" -table="third_payment_flow" -dir="." --style=goZero
