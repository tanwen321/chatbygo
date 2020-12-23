package api

import (
	// "fmt"
	"github.com/tanwen321/chatbygo/pkg/logging"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/tanwen321/chatbygo/models"
	"github.com/tanwen321/chatbygo/pkg/e"
	"github.com/tanwen321/chatbygo/pkg/util"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.PostForm("user")
	password := c.PostForm("passwd")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		userid, isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				ok := models.OnlineUser(userid, token)
				if ok {
					code = e.SUCCESS
				} else {
					code = e.ERROR_LOGINED_USER
				}
			}
		} else if userid == 0 {
			code = e.ERROR_USER_OR_PASS
		} else {
			code = e.ERROR_LOGINED_USER
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
