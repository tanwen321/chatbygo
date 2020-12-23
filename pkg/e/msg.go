package e

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_EXIST_USER:               "已存在该用户",
	ERROR_NOT_EXIST_USER:           "该用户不存在",
	ERROR_USER_PASS:                "用户密码错误",
	ERROR_USER_OR_PASS:             "用户或者密码错误",
	ERROR_USER_LOGINED:             "用户已经登陆",
	ERROR_NOT_EXIST_GROUP:          "该群组不存在",
	ERROR_EXIST_GROUP:              "该群组已经存在",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",
	ERROR_CREATE_USER:              "创建用户失败",
	ERROR_USER_NOT_ADMIN:           "用户权限不足",
	ERROR_CREATE_GROUP:             "创建群众失败",
	ERROR_USER_IN_GROUP:            "用户已经在该群组",
	ERROR_USER_ADD_FRIEND:          "该用户已经是好友",
	ERROR_LOGINED_USER:             "用户已经在其他地方登陆",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
