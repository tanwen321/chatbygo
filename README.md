# gchat
go语言的聊天服务器


### 1、运行环境
linux go环境或者mac go环境

### 2、依赖的go包
websocket：ws的基础包  
https://github.com/gorilla/websocket  
Gin：Golang 的一个微框架，性能极佳。  
https://github.com/gin-gonic/gin  
beego-validation：采用的 beego 的参数验证库。  
https://github.com/astaxie/beego/tree/master/validation  
gorm，对开发人员友好的 ORM 框架（数据库连接）。  
https://github.com/jinzhu/gorm  
com，一个小而美的工具包。  
https://github.com/Unknwon/com  

注意：提示缺少就使用  
go get -u giturl地址  

### 3、运行步骤
导入数据库和实例表  
mysql_db.sql是msql数据库的实例数据  
默认的用户名和密码查看sql文件  

修改配置文件（主要是数据库连接信息）  
文件路径conf/app.ini  

[@tanwendeMacBook-Air:G_work]$ cd chatbygo/  
[@tanwendeMacBook-Air:chatbygo]$ ./install   
finished 
[@tanwendeMacBook-Air:chatbygo]$ go run main.go  

### 4、API
添加用户（管理员权限）  
/api/v1/adduser?name=iphone15&pass=test123456&nikename=我是最新的&words=没有最新，只有更新&token=token信息  

添加群组（管理员权限）  
/api/v1/addgroup?name=做工的人&word=做工的人，做不完的工&token=token信息  

添加当前用户到群组  
/api/v1/addusertogroup?groupname=做工的人&token=token信息  

添加好友  
/api/v1/adduserfriend?friendname=admin&token=token信息  

#### 5、效果图

![image](https://github.com/tanwen321/chatbygo/blob/main/xiaoguo.png)

