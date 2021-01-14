package logger

import (
	"Ankr-gin-ERC721/conf"
	"log"
	"os"
	"testing"
)

func Test(t *testing.T)  {
	dir, _ := conf.GetAppPath()
	logPath := dir + `\log\log.txt`

	_, err := os.Create(logPath)
	if err != nil {
		log.Fatal("create log file failed, error: ", err)
	}
}
