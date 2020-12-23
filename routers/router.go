package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/tanwen321/chatbygo/middleware/jwt"
	"github.com/tanwen321/chatbygo/pkg/setting"
	"github.com/tanwen321/chatbygo/routers/api"
	"github.com/tanwen321/chatbygo/routers/api/v1"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	webhome := setting.WEB_HOME

	r.StaticFile("/favicon.ico", webhome+"/favicon.ico")
	r.StaticFile("/", webhome+"/index.html")
	r.StaticFile("/index", webhome+"/index.html")
	r.StaticFS("/js/", http.Dir(webhome+"/js"))
	r.StaticFS("/css/", http.Dir(webhome+"/css"))
	r.StaticFS("/images/", http.Dir(webhome+"/images"))
	r.StaticFS("/touxiang/", http.Dir(webhome+"/touxiang"))

	r.POST("/login.html", api.GetAuth)
	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//聊天服务
		apiv1.GET("/mmchat_socket", v1.WsHandler)
		//创建用户(管理员权限)
		apiv1.GET("/adduser", v1.AddUser)
		//创建用户群组(管理员权限)
		apiv1.GET("/addgroup", v1.AddGroup)
		//添加用户到群组
		apiv1.GET("/addusertogroup", v1.AddUserToGroup)
		//添加用户为好友
		apiv1.GET("/adduserfriend", v1.AddUserFriend)
	}

	return r
}
