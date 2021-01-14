package syncData

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/eth/parseLog"
	"Ankr-gin-ERC721/pkg/setting"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
	"time"
)

func TestSync721EventData(t *testing.T) {
	client, err := ethclient.Dial("")
	if err != nil {
		t.Error(err)
		return
	}

	address := "0xC7e5e9434f4a71e6dB978bD65B4D61D3593e5f27"
	//address := "0xC469E9Bd0276E6185C9A0B4E61fed8e0D0eD0185"
	//chainID := conf.MAIN_CHAINID
	contractAddress := common.HexToAddress(address)

	//	获取所有交易
	/*txes, err := interactContract.GetAllTXFromContract(address, chainID)
	if err != nil {
		t.Error(err)
		return
	}
	for i, tx := range txes {
		fmt.Printf("%d => %s\n", i, tx.Input)
	}
	var (
		startNumStr string
		startNum    int64
	)
	if len(txes) != 0 {
		startNumStr = txes[0].BlockNumber
		startNum, _ = strconv.ParseInt(startNumStr, 0, 0)
	}*/

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: nil,
		ToBlock:   nil,
	}
	logsOld, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		t.Error("ETHClient.FilterLogs error: ", err)
		return
	}
	for _, log := range logsOld {
		ok, from, to, tokenID := parseLog.ParseLog721(log)
		if ok {
			fmt.Printf("监听到事件：from %s to %s tokenID %d\n", from, to, tokenID)
		}
	}
}

func TestSubscribe(t *testing.T) {
HEAR:
	client, err := ethclient.Dial("wss://binance-sc-01.dccn.ankr.com/ws")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(client)
	startTime := time.Now()

	/*
		使用wss
		service := "wss://dex.binance.org/api/ws/0xC469E9Bd0276E6185C9A0B4E61fed8e0D0eD0185"
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		if err != nil {
			t.Error(err)
			return
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			t.Error(err)
			return
		}*/

	//address := "0xD8d638BE21B4101e1858DE84d4540Aed2d02674d"
	address := "0xC469E9Bd0276E6185C9A0B4E61fed8e0D0eD0185"
	contractAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		t.Error(err)
		return
	}
	go func() {
		for {
			id, _ := client.ChainID(context.Background())
			fmt.Println("chainID:", id)
			time.Sleep(40 * time.Second)
		}
	}()
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("sub.Err : ", err)
			endTime := time.Now()
			fmt.Println("时间间隔 => ", endTime.Sub(startTime).Seconds())
			client.Close()
			goto HEAR
		case vLog := <-logs:
			fmt.Println("合约事件日志 ===> ", vLog.Topics[0].String()) // pointer to event log
			//	解析事件日志
			isTransfer, from, to, tokenID := parseLog.ParseLog721(vLog)
			fmt.Printf("isok?=> %t 监听到事件：from %s to %s tokenID %d\n", isTransfer, from, to, tokenID)
		}
	}

}

func TestSync1155EventData(t *testing.T) {
	contractAddr := "0xC7e5e9434f4a71e6dB978bD65B4D61D3593e5f27"
	chainID := 1
	Sync1155EventData(contractAddr, nil, nil, chainID)
}

func TestSync721EventData2(t *testing.T) {
	str := "0x0000000000000000000000000000000000000000000000000000000000000001"
	newInt := big.NewInt(0)
	i, ok := newInt.SetString(str, 0)
	fmt.Println(i.Int64(), ok)
}

func TestSyncContractEventLog(t *testing.T) {
	contractAddr := "0xd1039feb22ecf8047e3532e317570dbfa93eb9b5"
	chainID := 56
	address := common.HexToAddress(contractAddr)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		FromBlock: big.NewInt(1363351),
		ToBlock:   big.NewInt(1368350),
	}
	//logger.Logger.Info().Msgf("start get logs from eth client, it will take a few seconds... address: %s chainID %d fromBLock %d toBLock %d", contractAddress, chainID, fromBlock.Int64(), toBlock.Int64())
	setting.ClientLock.Lock()
	logsOld, err := setting.ETHHTTPClients[chainID].FilterLogs(context.Background(), query)
	setting.ClientLock.Unlock()
	time := 0
	for err != nil && time < conf.RETRY_TIME {
		setting.ClientLock.Lock()
		logsOld, err = setting.ETHHTTPClients[chainID].FilterLogs(context.Background(), query)
		setting.ClientLock.Unlock()
		time++
	}
	if err != nil {
		t.Error(err)
		return
	}
	for _, log := range logsOld {
		_, from, to, tokenID := parseLog.ParseLog721(log)
		fmt.Println(from,to,tokenID)
	}
}
