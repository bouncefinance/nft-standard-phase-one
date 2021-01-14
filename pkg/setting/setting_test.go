package setting

import (
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	"Ankr-gin-ERC721/pkg/eth/parseLog"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strconv"
	"testing"
)

func TestLoadETHClient(t *testing.T) {
	/*base64Data, err := base64.URLEncoding.DecodeString("huxulong:@Ankr123")
	if err!=nil{
		t.Error(err)
	}

	base64Data
	*/
	client, err := ethclient.Dial("")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(client)

	address := "0xD8d638BE21B4101e1858DE84d4540Aed2d02674d"
	chainID := 4
	contractAddress := common.HexToAddress(address)

	//	获取所有交易
	txes, err := interactContract.GetAllTXFromContract(address, chainID)
	if err != nil {
		t.Error(err)
	}
	var (
		startNumStr string
		startNum    int64
	)
	if len(txes) != 0 {
		startNumStr = txes[0].BlockNumber
		startNum, _ = strconv.ParseInt(startNumStr, 0, 0)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: big.NewInt(startNum),
		ToBlock:   nil,
	}
	logsOld, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		t.Error("ETHClient.FilterLogs error: ", err)
	}
	for _, log := range logsOld {
		ok, from, to, tokenID := parseLog.ParseLog721(log)
		if ok {
			fmt.Printf("监听到事件：from %s to %s tokenID %d\n", from, to, tokenID)
		}
	}
}
