package setting

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/logger"
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-ini/ini"
	"log"
	"sync"
	"time"
)

var (
	Cfg            *ini.File
	configPath     = "conf/config.ini"
	RunMode        string
	HTTPPort       int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ETHClients     map[int]*ethclient.Client
	ETHHTTPClients map[int]*ethclient.Client
	ClientLock     sync.Mutex
	EtherscanURLS  map[int]string
	URLLock        sync.Mutex
)

const (
	MAIN_NET_URL    = ""
	ROPSTEN_NET_URL = ""
	KOVAN_NET_URL   = ""
	RINKEBY_NET_URL = ""
	GOERLI_NET_URL  = ""
	BSC_NET_URL     = ""

	MAIN_ETHERSCAN_URL    = "https://api.etherscan.io/api"
	ROPSTEN_ETHERSCAN_URL = "https://api-ropsten.etherscan.io/api"
	KOVAN_ETHERSCAN_URL   = "https://api-kovan.etherscan.io/api"
	RINKEBY_ETHERSCAN_URL = "https://api-rinkeby.etherscan.io/api"
	GOERLI_ETHERSCAN_URL  = "https://api-goerli.etherscan.io/api"
	BSC_ETHERSCAN_URL     = "https://api.bscscan.com/api"
)

func init() {
	var err error
	Cfg, err = ini.Load(conf.Dir + "/" + configPath)
	if err != nil {
		panic(err)
	}

	ETHClients = make(map[int]*ethclient.Client)
	ETHHTTPClients = make(map[int]*ethclient.Client)
	EtherscanURLS = make(map[int]string)
	LoadBase()
	LoadServer()
	LoadETHClient()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadETHClient() {
	configData, err := conf.ParseConfig()
	if err != nil {
		log.Fatalf("Fail to get ParseConfig error: %v", err)
	}

	ETHClients[conf.MAIN_CHAINID], err = ethclient.Dial(MAIN_NET_URL + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("wss client", "dial wss ETH main client success"+MAIN_NET_URL).Msg("")

	ETHClients[conf.RINKEBY_CHAINID], err = ethclient.Dial(RINKEBY_NET_URL + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("wss client", "dial wss ETH rinkeby client success"+RINKEBY_NET_URL).Msg("")

	ETHClients[conf.ROPSTEN_CHAINID], err = ethclient.Dial(ROPSTEN_NET_URL + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("wss client", "dial wss ETH ropsten client success"+ROPSTEN_NET_URL).Msg("")

	ETHClients[conf.KOVAN_CHAINID], err = ethclient.Dial(KOVAN_NET_URL + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("wss client", "dial wss ETH kovan client success"+KOVAN_NET_URL).Msg("")

	ETHClients[conf.GOERLI_CHAINID], err = ethclient.Dial(GOERLI_NET_URL + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("wss client", "dial wss ETH goerli client success"+GOERLI_NET_URL).Msg("")

	ETHClients[conf.BSC_CHAINID], err = ethclient.Dial(BSC_NET_URL)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("wss client", "dial wss BSC main client success"+BSC_NET_URL).Msg("")

	ETHHTTPClients[conf.MAIN_CHAINID], err = ethclient.Dial(conf.MAIN_NET_URL + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("https client", "dial https ETH main client success"+conf.MAIN_NET_URL).Msg("")

	ETHHTTPClients[conf.RINKEBY_CHAINID], err = ethclient.Dial(conf.RINKEBY_NET_URL + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("https client", "dial https ETH main client success"+conf.RINKEBY_NET_URL).Msg("")

	ETHHTTPClients[conf.ROPSTEN_CHAINID], err = ethclient.Dial(conf.ROPSTEN_NET_URL + "/" + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("https client", "dial https ETH main client success"+conf.ROPSTEN_NET_URL).Msg("")

	ETHHTTPClients[conf.KOVAN_CHAINID], err = ethclient.Dial(conf.KOVAN_NET_URL + "/" + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("https client", "dial https ETH main client success"+conf.KOVAN_NET_URL).Msg("")

	ETHHTTPClients[conf.GOERLI_CHAINID], err = ethclient.Dial(conf.GOERLI_NET_URL + "/" + configData.ProjectID)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("https client", "dial https ETH main client success"+conf.GOERLI_NET_URL).Msg("")

	ETHHTTPClients[conf.BSC_CHAINID], err = ethclient.Dial(conf.BSC_NET_URL)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Str("https client", "dial https BSC main client success"+conf.BSC_NET_URL).Msg("")

	go func() {
		ticker := time.NewTicker(40 * time.Second)
		for {
			select {
			case <- ticker.C:
				_, _ = ETHClients[conf.BSC_CHAINID].ChainID(context.Background())
			}
		}
	}()

	URLLock.Lock()
	EtherscanURLS[conf.MAIN_CHAINID] = MAIN_ETHERSCAN_URL
	EtherscanURLS[conf.RINKEBY_CHAINID] = RINKEBY_ETHERSCAN_URL
	EtherscanURLS[conf.ROPSTEN_CHAINID] = ROPSTEN_ETHERSCAN_URL
	EtherscanURLS[conf.KOVAN_CHAINID] = KOVAN_ETHERSCAN_URL
	EtherscanURLS[conf.GOERLI_CHAINID] = GOERLI_ETHERSCAN_URL
	EtherscanURLS[conf.BSC_CHAINID] = BSC_ETHERSCAN_URL
	URLLock.Unlock()
}
