package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

const (
	LOG_PATH = "./log.txt"
)

type GetErc struct {
	UserAddr     string `json:"userAddr"`
	ContractAddr string `json:"contractAddr"`
	ChainID      int    `json:"chainID"`
	Message      string `json:"message"`
}

type GetAllNFT struct {
	UserAddr string `json:"userAddr"`
	ChainID  int    `json:"chainID"`
	Message  string `json:"message"`
}

type Message struct {
	Message string `json:"message"`
}

func main() {
	file, err := os.OpenFile(LOG_PATH, os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("openFile failed error:", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	UserRequest(r,"0x49963648b69a9d318e7ca9f36d2b890d0fef9b5f")
}

func OpenLogFile() *bufio.Reader {
	file, err := os.OpenFile(LOG_PATH, os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("openFile failed error:", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	return r
}

func UserRequest(reader *bufio.Reader,subStr string) {
	requestLog, err := os.Create("./log_deal.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer requestLog.Close()
	//msg := Message{}
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}
		/*err = json.Unmarshal(buf, &msg)
		if err!=nil{
			fmt.Println("json.Unmarshal err: ",err)
		}*/
		if strings.Contains(string(buf), subStr) {
			_, err := requestLog.Write(buf)
			if err != nil {
				log.Fatal("write file failed", err)
			}
		}
	}
}
