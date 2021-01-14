package routers

import (
	"Ankr-gin-ERC721/middleware"
	"Ankr-gin-ERC721/pkg/setting"
	"Ankr-gin-ERC721/routers/api"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.TlsHandler())

	gin.SetMode(setting.RunMode)

	r.GET("/erc721",api.GetERC721)
	r.GET("/erc1155",api.GetERC1155)
	r.GET("/nft",api.GetAllNFT)

	assetsGroup:=r.Group("/assets")
	{
		assetsGroup.GET("/erc721",api.GetMetadata721)
		assetsGroup.GET("/erc1155",api.GetMetadata1155)
	}

	return r
}

