package nfts

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/setting"
	"Ankr-gin-ERC721/pkg/util"
	"Ankr-gin-ERC721/routers/api/packaging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
	"sync"
	"time"
)

func All721NFTAssets(address string, chainID int) error {
	time := 0
	contracts := make([]string, 0)
	txReceipts, err := interactContract.GetAll721TX(address, chainID)
	for err != nil && time < conf.RETRY_TIME {
		txReceipts, err = interactContract.GetAll721TX(address, chainID)
		time++
	}
	if err != nil {
		return errors.Wrap(err, "GetAll721TX error")
	}
	for _, txReceipt := range txReceipts {
		has := false
		for _, contractAddr := range contracts {
			if contractAddr == txReceipt.ContractAddress {
				has = true
				break
			}
		}
		if has {
			continue
		}
		contracts = append(contracts, txReceipt.ContractAddress)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(contracts))
	for _, contract := range contracts {
		contract = util.StrToLow(contract)
		go func(contract string) {
			packaging.Launch721(contract, chainID)
			defer wg.Done()
		}(contract)
	}
	wg.Wait()
	logger.Logger.Info().Str("userAddress",address).Int("chainID",chainID).Msg("All721NFT同步结束")
	return nil
}

func All1155NFTAssets(address string, chainID int, startNum, endNum int64) {
	contractAddrs := make([]string, 0)
	networkID, err := setting.ETHClients[chainID].NetworkID(context.Background())
	if err != nil {
		fmt.Println("client NetworkID error: ", err)
		return
	}
	var i uint64
	wg := sync.WaitGroup{}
	wg.Add(int(endNum - startNum + 1))
	for i = uint64(startNum); i <= uint64(endNum); {
		for j := 0; j < 5; j++ {
			i++
			if i > uint64(endNum) {
				break
			}

			go func(i uint64) {
				defer wg.Done()
				block, err := setting.ETHHTTPClients[chainID].BlockByNumber(context.Background(), big.NewInt(int64(i)))
				if err != nil {
					fmt.Printf("client BlockByNumber error: %v point:%v \n", err, setting.ETHHTTPClients[chainID])
					return
				}
				transactions := block.Transactions()
				fmt.Printf("blockNum: %d endNum: %d txLength: %d \n", i+1, endNum, transactions.Len())
				for _, transaction := range transactions {
					var to string
					toTem := transaction.To()
					if toTem != nil {
						to = toTem.Hex()
					}
					message, err := transaction.AsMessage(types.NewEIP155Signer(networkID))
					if err != nil {
						fmt.Println("transaction.AsMessage error: ", err)
						continue
					}
					from := message.From().Hex()

					if from == address || to == address {
						fmt.Printf("from : %s to : %s isOk: %t\n", from, to, from == address || to == address)
						receipt, err := setting.ETHClients[chainID].TransactionReceipt(context.Background(), transaction.Hash())
						if err != nil {
							fmt.Println("client TransactionReceipt error: ", err)
							continue
						}
						data, err := interactContract.SupportInterface(receipt.ContractAddress.Hex(), interactContract.ERC1155_INTERFACE_ID, chainID)
						if err != nil {
							fmt.Println("interactContract.SupportInterface error: ", err)
							continue
						}
						returnData := interactContract.ReturnData{}
						err = json.Unmarshal(data, &returnData)
						if err != nil {
							fmt.Println("All1155NFTAssets json.Unmarshal error: ", err)
							continue
						}
						iTem, err := strconv.ParseInt(returnData.Result, 0, 0)
						if err != nil {
							fmt.Println("All1155NFTAssets strconv.ParseInt error: ", err)
							continue
						}
						if iTem == interactContract.ACTIVELY_RESPONSE {
							contractAddrs = append(contractAddrs, receipt.ContractAddress.Hex())
						}
					}
				}
			}(i)
		}
		time.Sleep(800 * time.Millisecond)
	}
	wg.Wait()
	wg.Add(len(contractAddrs))
	for _, addr := range contractAddrs {
		go func(addr string) {
			packaging.Launch1155(addr, chainID)
			defer wg.Done()
		}(addr)
	}
	wg.Wait()
}
