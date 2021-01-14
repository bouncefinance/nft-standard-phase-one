package conf

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	MAIN_NET_URL    = ""
	ROPSTEN_NET_URL = ""
	KOVAN_NET_URL   = ""
	RINKEBY_NET_URL = ""
	GOERLI_NET_URL  = ""
	BSC_NET_URL     = ""
	BSC_TEST_URL    = ""

	MAIN_CHAINID     = 1
	ROPSTEN_CHAINID  = 2
	KOVAN_CHAINID    = 3
	RINKEBY_CHAINID  = 4
	GOERLI_CHAINID   = 5
	BSC_CHAINID      = 56
	BSC_TEST_CHAINID = 97

	MONGO_UTL = ""
)

type ConfigData struct {
	Mnemonic    string `json:"mnemonic"`
	ProjectID   string `json:"project_id"`
	ApiKeyToken string `json:"api_key_token"`
}

var (
	Dir         string
	ConfigPath  string
	ConfigMsg   ConfigData
	ETHHttpURLs map[int]string
	err         error
)

func GetAppPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	p, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	index := strings.LastIndex(p, string(os.PathSeparator))
	return p[:index], nil
}

func ParseConfig() (config ConfigData, err error) {
	file, err := os.OpenFile(ConfigPath, os.O_RDONLY, 0666)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return
	}
	return
}
func init() {
	Dir, _ = GetAppPath()
	ConfigPath = Dir + "/conf/config.json"
	ConfigMsg, err = ParseConfig()
	if err != nil {
		log.Fatal("ParseConfig error: ", err)
	}

	ETHHttpURLs = make(map[int]string)
	ETHHttpURLs[MAIN_CHAINID] = MAIN_NET_URL
	ETHHttpURLs[ROPSTEN_CHAINID] = ROPSTEN_NET_URL
	ETHHttpURLs[KOVAN_CHAINID] = KOVAN_NET_URL
	ETHHttpURLs[RINKEBY_CHAINID] = RINKEBY_NET_URL
	ETHHttpURLs[GOERLI_CHAINID] = GOERLI_NET_URL
	ETHHttpURLs[BSC_CHAINID] = BSC_NET_URL
	ETHHttpURLs[BSC_TEST_CHAINID] = BSC_TEST_URL
}
