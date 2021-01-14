package interactContract

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/setting"
	"Ankr-gin-ERC721/pkg/util"
	"encoding/json"
	"github.com/pkg/errors"
)

type TXReceipt struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

type EtherscanTXData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  []TXReceipt `json:"result"`
}

var (
	NotFoundError = errors.New("There are not transactions.")
)

const RETRY_TIME = 5

func GetAllNormalTX(address string, chainID int) ([]TXReceipt, error) {
	time := 0
	params := make(map[string]string)
	params["module"] = "account"
	params["action"] = "txlist"
	params["address"] = address
	params["apikey"] = conf.ConfigMsg.ApiKeyToken
	data, err := util.PostUrl(setting.EtherscanURLS[chainID], params, nil, nil)
	for err != nil && time < RETRY_TIME {
		data, err = util.PostUrl(setting.EtherscanURLS[chainID], params, nil, nil)
		time++
	}
	if err != nil {
		return nil, errors.Wrap(err,"util.PostUrl error")
	}

	etherscanTXData := EtherscanTXData{}
	err = json.Unmarshal(data, &etherscanTXData)
	if err != nil {
		return nil, errors.Wrap(err,"json.Unmarshal error")
	}
	if etherscanTXData.Status == "0" {
		return nil, NotFoundError
	}
	return etherscanTXData.Result, nil
}

func GetNormalTXWithBlockNum(address string, chainID int,fromBlock ,toBlock string) ([]TXReceipt, error) {
	time := 0
	params := make(map[string]string)
	params["module"] = "account"
	params["action"] = "txlist"
	params["address"] = address
	params["startblock"] = fromBlock
	params["endblock"] = toBlock
	params["apikey"] = conf.ConfigMsg.ApiKeyToken
	data, err := util.PostUrl(setting.EtherscanURLS[chainID], params, nil, nil)
	for err != nil && time < RETRY_TIME {
		data, err = util.PostUrl(setting.EtherscanURLS[chainID], params, nil, nil)
		time++
	}
	if err != nil {
		return nil, errors.Wrap(err,"util.PostUrl error")
	}

	etherscanTXData := EtherscanTXData{}
	err = json.Unmarshal(data, &etherscanTXData)
	if err != nil {
		return nil, errors.Wrap(err,"json.Unmarshal error")
	}
	if etherscanTXData.Status == "0" {
		return nil, NotFoundError
	}
	return etherscanTXData.Result, nil
}

func GetAll721TX(address string, chainID int) ([]TXReceipt, error) {
	params := make(map[string]string)
	params["module"] = "account"
	params["action"] = "tokennfttx"
	params["address"] = address
	params["apikey"] = conf.ConfigMsg.ApiKeyToken
	data, err := util.PostUrl(setting.EtherscanURLS[chainID], params, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err,"util.PostUrl error")
	}

	etherscanTXData := EtherscanTXData{}
	err = json.Unmarshal(data, &etherscanTXData)
	if err != nil {
		return nil, errors.Wrap(err,"json.Unmarshal error")
	}
	if etherscanTXData.Status == "0" {
		return nil, errors.New(etherscanTXData.Message)
	}
	return etherscanTXData.Result, nil
}
