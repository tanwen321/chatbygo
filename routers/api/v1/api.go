package v1

import (
	"github.com/tanwen321/chatbygo/pkg/logging"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/tanwen321/chatbygo/models"
	"github.com/tanwen321/chatbygo/pkg/e"
)

func AddUser(c *gin.Context) {

	token := c.Query("token")
	id, ok1 := models.GetUserBytoken(token)
	ok2 := models.CheckAdmin(id)

	code := e.INVALID_PARAMS
	msg := ""
	if ok1 && ok2 {
		name := c.Query("name")
		pass := c.Query("pass")
		nikename := c.Query("nikename")
		words := c.DefaultQuery("words", "这家伙很懒，什么都没留下")
		touxiang := randimage()
		rool := com.StrTo(c.DefaultQuery("rool", "10")).MustInt()
		status := com.StrTo(c.DefaultQuery("status", "1")).MustInt()

		valid := validation.Validation{}
		valid.Required(name, "name").Message("用户名不能为空")
		valid.MaxSize(name, 40, "name").Message("用户名最长为40字符")
		valid.MinSize(name, 4, "name").Message("用户名最短为4字符")
		valid.Required(pass, "pass").Message("密码不能为空")
		valid.MaxSize(pass, 40, "pass").Message("密码最长为40字符")
		valid.MinSize(pass, 6, "pass").Message("密码最短为6字符")
		valid.Required(nikename, "nikename").Message("昵称不能为空")
		valid.MaxSize(nikename, 90, "nikename").Message("昵称最长为90字符")

		if !valid.HasErrors() {
			uid := models.ExistUserByName(name)
			if uid == 0 {
				data := make(map[string]interface{})
				data["name"] = name
				data["pass"] = pass
				data["nikename"] = nikename
				data["words"] = words
				data["rool"] = rool
				data["touxiang"] = touxiang
				data["status"] = status
				models.AddUser(data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_EXIST_USER
			}
		} else {
			for _, err := range valid.Errors {
				logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
			}
			code = 0
			msg = valid.Errors[0].Message
		}
	} else {
		code = e.ERROR_USER_NOT_ADMIN
	}

	if code == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  msg,
			"data": make(map[string]interface{}),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}
}

func AddGroup(c *gin.Context) {

	token := c.Query("token")
	id, ok1 := models.GetUserBytoken(token)
	ok2 := models.CheckAdmin(id)

	code := e.INVALID_PARAMS
	msg := ""
	if ok1 && ok2 {
		name := c.Query("name")
		words := c.DefaultQuery("words", "你来或者不来，群就在这里")
		touxiang := randimage()
		status := com.StrTo(c.DefaultQuery("status", "1")).MustInt()

		valid := validation.Validation{}
		valid.Required(name, "name").Message("群组名不能为空")
		valid.MaxSize(name, 40, "name").Message("群组名最长为40字符")
		valid.MinSize(name, 4, "name").Message("群组名最短为4字符")

		if !valid.HasErrors() {
			gid := models.ExistGroupByName(name)
			if gid == 0 {
				data := make(map[string]interface{})
				data["name"] = name
				data["words"] = words
				data["touxiang"] = touxiang
				data["status"] = status
				models.AddGroup(data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_EXIST_GROUP
			}

		} else {
			for _, err := range valid.Errors {
				logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
			}
			code = 0
			msg = valid.Errors[0].Message
		}
	} else {
		code = e.ERROR_USER_NOT_ADMIN
	}

	if code == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  msg,
			"data": make(map[string]interface{}),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}

}

func AddUserToGroup(c *gin.Context) {

	token := c.Query("token")
	id, ok := models.GetUserBytoken(token)

	code := e.INVALID_PARAMS
	msg := ""
	if ok {
		groupname := c.Query("groupname")

		valid := validation.Validation{}
		valid.Required(groupname, "groupname").Message("群组名不能为空")
		valid.MaxSize(groupname, 40, "groupname").Message("群组名最长为40字符")
		valid.MinSize(groupname, 4, "groupname").Message("群组名最短为4字符")

		if !valid.HasErrors() {
			gid := models.ExistGroupByName(groupname)
			if gid > 0 {
				if models.AddUserToGroup(id, gid) {
					code = e.SUCCESS
				} else {
					code = e.ERROR_USER_IN_GROUP
				}
			} else {
				code = e.ERROR_NOT_EXIST_GROUP
			}

		} else {
			for _, err := range valid.Errors {
				logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
			}
			code = 0
			msg = valid.Errors[0].Message
		}
	} else {
		code = e.ERROR_NOT_EXIST_USER
	}

	if code == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  msg,
			"data": make(map[string]interface{}),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}
}

func AddUserFriend(c *gin.Context) {

	token := c.Query("token")
	id, ok := models.GetUserBytoken(token)

	code := e.INVALID_PARAMS
	msg := ""
	if ok {
		friendname := c.Query("friendname")

		valid := validation.Validation{}
		valid.Required(friendname, "friendname").Message("好友名不能为空")
		valid.MaxSize(friendname, 40, "friendname").Message("好友名最长为40字符")
		valid.MinSize(friendname, 4, "friendname").Message("好友名最短为4字符")

		if !valid.HasErrors() {
			uid := models.ExistUserByName(friendname)
			if uid > 0 {
				if models.AddUserToFriend(id, uid) {
					code = e.SUCCESS
				} else {
					code = e.ERROR_USER_ADD_FRIEND
				}
			} else {
				code = e.ERROR_NOT_EXIST_USER
			}

		} else {
			for _, err := range valid.Errors {
				logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
			}
			code = 0
			msg = valid.Errors[0].Message
		}
	} else {
		code = e.ERROR_NOT_EXIST_USER
	}

	if code == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  msg,
			"data": make(map[string]interface{}),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}
}
