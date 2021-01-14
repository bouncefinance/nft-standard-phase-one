package msg

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",



	ERROR_CHECK_TOKEN:         "Token验证失败",
	ERROR_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_CHECK_TOKEN_NULL:    "Token为空",
	ERROR_ACCESS_TOKEN_NULL:   "访问小程序服务接口的token为空",
}


func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

