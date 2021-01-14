package subscribe

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/mongoData"
	"Ankr-gin-ERC721/pkg/eth/parseLog"
	"Ankr-gin-ERC721/pkg/eth/syncData"
	"Ankr-gin-ERC721/pkg/eventLoop"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/setting"
	"Ankr-gin-ERC721/runtime"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"math/big"
)

type Data721 struct {
	ContractAddr string
	ChainID      int
}

type Data1155 struct {
	ContractAddr string
	FromBlock    *big.Int
	ChainID      int
}

func Subscribe721(data *eventLoop.EventData) bool {
	subscribeData := data.Data.(Data721)

	contractAddress := common.HexToAddress(subscribeData.ContractAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	logs := make(chan types.Log)
	setting.ClientLock.Lock()
	sub, err := setting.ETHClients[subscribeData.ChainID].SubscribeFilterLogs(context.Background(), query, logs)
	setting.ClientLock.Unlock()
	time := 0
	for err != nil && time < conf.RETRY_TIME {
		setting.ClientLock.Lock()
		sub, err = setting.ETHClients[subscribeData.ChainID].SubscribeFilterLogs(context.Background(), query, logs)
		setting.ClientLock.Unlock()
		time++
	}
	if err != nil {
		logger.Logger.Error().Str("contractAddress", subscribeData.ContractAddr).Int("chainID", subscribeData.ChainID).Msgf("Subscribe721 client SubscribeFilterLogs error: %s", err)
		return false
	}

	logger.Logger.Info().Str("contractAddress", subscribeData.ContractAddr).Int("chainID", subscribeData.ChainID).Msg("开始订阅721合约...")
	for {
		select {
		case err := <-sub.Err():
			logger.Logger.Error().Str("contractAddress", subscribeData.ContractAddr).Int("chainID", subscribeData.ChainID).Msgf("Subscribe721 sub.Err error: %s", err)
			if subscribeData.ChainID == conf.BSC_CHAINID {
				setting.ETHClients[subscribeData.ChainID], _ = ethclient.Dial(conf.BSC_NET_URL)
			}
		case vLog := <-logs:
			isTransfer, from, to, tokenID := parseLog.ParseLog721(vLog)
			if isTransfer {
				logger.Logger.Info().Str("contractAddress", vLog.Address.String()).Int("chainID", subscribeData.ChainID).
					Str("from",from).Str("to",to).Int("tokenID",tokenID).
					Msg("监听到721合约transfer事件日志")
				go func() {
					err := syncData.UpdateNFT721Data(tokenID, subscribeData.ContractAddr, from, to, subscribeData.ChainID)
					if err != nil {
						logger.Logger.Error().Str("contractAddress", subscribeData.ContractAddr).Int("chainID", subscribeData.ChainID).Int("tokenID", tokenID).Str("from", from).Str("to", to).Msgf("Subscribe721 UpdateNFT721Data error:%s", err)
						runtime.EventLoop.Off(runtime.SUBSCRIBE_CONTRACT_721+subscribeData.ContractAddr, Subscribe721)
						return
					}
					filter := bson.D{{"ContractAddress", contractAddress}}
					updateData := bson.D{
						{"$max", bson.D{
							{"BlockNum", int(vLog.BlockNumber)},
						}},
					}
					mongoData.MongoDB.UpdateOne(mongoData.DATABASE, mongoData.COLLECTION, filter, updateData)
					if err != nil {
						logger.Logger.Error().Str("contractAddress", vLog.Address.String()).Int("ChainID", subscribeData.ChainID).Uint64("LatestBlockNum", vLog.BlockNumber).Msgf("Subscribe721 mongo UpdateOne error: %s", err)
					}
				}()
			}
		}
	}
	return true
}

func Subscribe1155(data *eventLoop.EventData) bool {
	data1155 := data.Data.(Data1155)

	contractAddress := common.HexToAddress(data1155.ContractAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: data1155.FromBlock,
	}
	logs := make(chan types.Log)

	logger.Logger.Info().Str("contractAddress", data1155.ContractAddr).Int("chainID", data1155.ChainID).Int64("start block num", data1155.FromBlock.Int64()).Msg("开始订阅1155合约...")
	setting.ClientLock.Lock()
	sub, err := setting.ETHClients[data1155.ChainID].SubscribeFilterLogs(context.Background(), query, logs)
	setting.ClientLock.Unlock()
	time := 0
	for err != nil && time < conf.RETRY_TIME {
		setting.ClientLock.Lock()
		sub, err = setting.ETHClients[data1155.ChainID].SubscribeFilterLogs(context.Background(), query, logs)
		setting.ClientLock.Unlock()
		time++
	}
	if err != nil {
		logger.Logger.Error().Str("contractAddress", data1155.ContractAddr).Int("chainID", data1155.ChainID).Msgf("Subscribe1155 client SubscribeFilterLogs error: %s", err)
		return false
	}
	for {
		select {
		case err := <-sub.Err():
			logger.Logger.Error().Str("contractAddress", data1155.ContractAddr).Int("chainID", data1155.ChainID).Msgf("Subscribe1155 sub.Err error: %s", err)
		case vLog := <-logs:
			flag, from, to, tokenIDs, amounts, uri, tokenID := parseLog.ParseLog1155(vLog)
			if flag != parseLog.MISSED {
				logger.Logger.Info().Str("contractAddress", vLog.Address.String()).Int("chainID", data1155.ChainID).
					Int("flag", flag).Str("from", from).Str("to", to).Uints("ids", tokenIDs).
					Interface("amounts", amounts).Msg("监听到1155合约事件日志，解析完毕")
				go func() {
					if flag == parseLog.TRANSFERSINGLE || flag == parseLog.TRANSFERBATCH {
						err := syncData.UpdateSubscribe1155Data(tokenIDs, amounts, data1155.ContractAddr, from, to, data1155.ChainID)
						if err != nil {
							logger.Logger.Error().Str("contractAddress", data1155.ContractAddr).Int("chainID", data1155.ChainID).Uints("tokenIDs", tokenIDs).Str("from", from).Str("to", to).Msgf("Subscribe1155 UpdateSubscribe1155Data error:%s", err)
							runtime.EventLoop.On(runtime.SUBSCRIBE_CONTRACT_1155+data1155.ContractAddr, Subscribe1155)
							return
						}
					}
					if flag == parseLog.URI {
						err := syncData.Update1155URI(data1155.ContractAddr, tokenID, uri)
						if err != nil {
							logger.Logger.Error().Str("contractAddress", data1155.ContractAddr).Int("chainID", data1155.ChainID).Int("tokenID", tokenID).Msgf("Subscribe1155 Update1155URI error:%s", err)
							return
						}
					}
					filter := bson.D{{"ContractAddress", data1155.ContractAddr}}
					updateData := bson.D{
						{"$max", bson.D{
							{"BlockNum", int(vLog.BlockNumber)},
						}},
					}
					err = mongoData.MongoDB.UpdateOne(mongoData.DATABASE, mongoData.COLLECTION,filter,updateData)
					if err != nil {
						logger.Logger.Error().Str("contractAddress", vLog.Address.String()).Int("ChainID", data1155.ChainID).Uint64("LatestBlockNum", vLog.BlockNumber).Msgf("Subscribe1155 mongo UpdateOne error: %s", err)
					}
				}()
			}
		}
	}
	return true
}
