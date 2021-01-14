package main

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/setting"
	"Ankr-gin-ERC721/routers"
	"fmt"
	"net/http"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}


	logger.Logger.Info().Int("port",setting.HTTPPort).Msgf("service run at %d",setting.HTTPPort)
	router.RunTLS(":443",conf.Dir + "/" + "cert/4618025_nftview.bounce.finance.pem",conf.Dir + "/" + "cert/4618025_nftview.bounce.finance.key")
	s.ListenAndServe()
}
