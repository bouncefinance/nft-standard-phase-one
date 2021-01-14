package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGetUrl(t *testing.T) {
	sendRequest()
}

type GetAllNFT struct {
	UserAddr string `json:"userAddress"`
	ChainID  int    `json:"chainID"`
	Message  string `json:"message"`
}

func sendRequest() {
	url := "https://nftview.bounce.finance/nft"

	file, err := os.OpenFile("log\\log_request.txt", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("openFile failed error:", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	wg := sync.WaitGroup{}
	for {
		time.Sleep(50 * time.Millisecond)
		buf, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}

		param := GetAllNFT{}

		if strings.Contains(string(buf), "GetAllNFT") {
			err = json.Unmarshal(buf, &param)
			if err != nil {
				log.Fatal("json.Unmarshal err: ", err)
			}
			go func(param GetAllNFT) {
				wg.Add(1)
				defer wg.Done()

				data, err := GetUrl(url, map[string]string{"address": param.UserAddr, "chain_id": strconv.Itoa(param.ChainID)})
				if err != nil {
					log.Fatal("GetUrl error: ", err)
				}
				fmt.Println(string(data))
			}(param)
		}
	}
	wg.Wait()

}
