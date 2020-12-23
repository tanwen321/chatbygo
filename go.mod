module github.com/tanwen321/chatbygo

go 1.14

require (
	github.com/EDDYCJY/go-gin-example v0.0.0-20200505102242-63963976dee0
	github.com/astaxie/beego v1.12.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.62.0
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/ugorji/go v1.2.1 // indirect
	github.com/unknwon/com v1.0.1
	golang.org/x/crypto v0.0.0-20201217014255-9d1352758620 // indirect
	golang.org/x/sys v0.0.0-20201214210602-f9fddec55a1e // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/tanwen321/chatbygo/conf => /Users/tanwen/Desktop/G_work/chatbygo/pkg/conf
	github.com/tanwen321/chatbygo/middleware => /Users/tanwen/Desktop/G_work/chatbygo/middleware
	github.com/tanwen321/chatbygo/models => /Users/tanwen/Desktop/G_work/chatbygo/models
	github.com/tanwen321/chatbygo/pkg/e => /Users/tanwen/Desktop/G_work/chatbygo/pkg/e
	github.com/tanwen321/chatbygo/pkg/logging => /Users/tanwen/Desktop/G_work/chatbygo/pkg/logging
	github.com/tanwen321/chatbygo/pkg/setting => /Users/tanwen/Desktop/G_work/chatbygo/pkg/setting
	github.com/tanwen321/chatbygo/pkg/util => /Users/tanwen/Desktop/G_work/chatbygo/pkg/util
	github.com/tanwen321/chatbygo/routers => /Users/tanwen/Desktop/G_work/chatbygo/routers
	github.com/tanwen321/chatbygo/routers/api => /Users/tanwen/Desktop/G_work/chatbygo/routers/api
)
