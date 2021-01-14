package logger

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/util"
	"fmt"
	"github.com/rs/zerolog"
	"log"
	"os"
	"time"
)

var Logger zerolog.Logger

func init() {
	logTime := time.Now().Unix()
	dir, _ := conf.GetAppPath()
	logPath := dir + fmt.Sprintf("/log/log%d.txt", logTime)
	file := &os.File{}
	var err error
	if util.FileNotExist(logPath) {
		file, err = os.Create(logPath)
		if err != nil {
			log.Fatal("create log file failed, error: ", err)
		}
	}else {
		file, err = os.OpenFile(logPath,os.O_RDWR,0666)
		if err != nil {
			log.Fatal("open log file failed, error: ", err)
		}
	}

	Logger = zerolog.New(file).With().Timestamp().Logger()
	Logger.Level(zerolog.InfoLevel)
}
