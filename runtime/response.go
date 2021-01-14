package runtime

import (
	"Ankr-gin-ERC721/pkg/msg"
	"github.com/gin-gonic/gin"
)

type Context struct {
	C *gin.Context
}

func (c *Context) Response(httpCode, errCode int, data interface{}) {
	c.C.Header("Access-Control-Allow-Origin","*")
	c.C.Header("Access-Control-Allow-Methods","*")

	c.C.JSON(httpCode, gin.H{
		"codeStatus": errCode,
		"msg":  msg.GetMsg(errCode),
		"data": data,
	})
	return
}

func (c *Context) ResponseMetaData(httpCode, errCode int, data interface{}) {
	c.C.Header("Access-Control-Allow-Origin","*")
	c.C.Header("Access-Control-Allow-Methods","*")

	c.C.JSON(httpCode, gin.H{
		"assets": data,
	})
	return
}
