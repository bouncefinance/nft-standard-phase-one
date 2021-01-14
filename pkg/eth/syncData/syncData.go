package syncData

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/mongoData"
	"Ankr-gin-ERC721/pkg/cache"
	"Ankr-gin-ERC721/pkg/eth/parseLog"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/setting"
	"Ankr-gin-ERC721/pkg/util"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"math/big"
	"sync"
)

const (
	DIVID_CONTAIN = 5000
)

func Sync1155EventData(contractAddr string, fromBlock *big.Int, toBlock *big.Int, chainID int) (blockNum uint64, err error) {
	contractAddress := common.HexToAddress(contractAddr)
	time := DivideTime(fromBlock, toBlock)
	eventSigns := []common.Hash{parseLog.SignBatch, parseLog.SignSingle, parseLog.SignURI}

	logs := commonArea(contractAddress, time, toBlock, fromBlock, chainID, eventSigns)
	filter := bson.D{{"ContractAddress", contractAddr}, {"ChainID", chainID}}

	for i, log := range logs {
		flag, from, to, ids, values, uri, tokenID := parseLog.ParseLog1155(log)
		if flag != parseLog.MISSED {
			logger.Logger.Info().Str("contractAddress", contractAddr).Int("chainID", chainID).Int("flag", flag).Str("from", from).Str("to", to).Uints("ids", ids).Interface("amounts", values).Msg("解析出1155合约日志")
			if flag == parseLog.TRANSFERSINGLE || flag == parseLog.TRANSFERBATCH {
				err = UpdateSubscribe1155Data(ids, values, contractAddr, from, to, chainID)
				if err != nil {
					logger.Logger.Error().Str("contractAddress", contractAddr).Int("chainID", chainID).Uints("ids", ids).Interface("amounts", values).Str("from", from).Str("to", to).Msgf("Sync1155EventData UpdateSubscribe1155Data error:%s", err)
					if i == 0 {
						blockNum = 0
					} else {
						blockNum = logs[i-1].BlockNumber
					}
					err = errors.Wrap(err, "UpdateNFT721Data error")
					return
				}
				updateData := bson.D{
					{"$max", bson.D{
						{"BlockNum", int(log.BlockNumber)},
					}},
				}
				err := mongoData.MongoDB.UpdateOne(mongoData.DATABASE, mongoData.COLLECTION, filter, updateData)
				if err != nil {
					logger.Logger.Error().Str("contractAddress", contractAddr).Int("ChainID", chainID).Uint64("LatestBlockNum", log.BlockNumber).Msgf("Launch1155 error: %s", err)
				}
			}
			if flag == parseLog.URI {
				err := Update1155URI(contractAddr, tokenID, uri)
				if err != nil {
					logger.Logger.Error().Str("contractAddress", contractAddr).Int("chainID", chainID).Int("tokenID", tokenID).Msgf("Sync1155EventData Update1155URI error:%s", err)
					continue
				}

			}
		}
	}
	logger.Logger.Info().Str("contractAddress", contractAddr).Int("chainID", chainID).
		Int64("from", fromBlock.Int64()).Int64("to", toBlock.Int64()).Msg("此1155合约数据同步完毕")
	cache.SyncSignMap.Delete(contractAddr)
	if len(logs) == 0 {
		return toBlock.Uint64(), nil
	} else {
		return logs[len(logs)-1].BlockNumber, nil
	}
}

func Sync721EventData(contractAddr string, fromBlock *big.Int, toBlock *big.Int, chainID int) (blockNum uint64, err error) {
	contractAddress := common.HexToAddress(contractAddr)
	time := DivideTime(fromBlock, toBlock)
	eventSigns := []common.Hash{parseLog.SignTransfer721}

	logs := commonArea(contractAddress, time, toBlock, fromBlock, chainID, eventSigns)
	filter := bson.D{{"ContractAddress", contractAddr}, {"ChainID", chainID}}

	for i, log := range logs {
		ok, from, to, tokenID := parseLog.ParseLog721(log)
		if ok {
			logger.Logger.Info().Str("contractAddress", contractAddr).Int("chainID", chainID).
				Str("from", from).Str("to", to).Int("tokenID", tokenID).Msg("同步数据，解析出721合约日志")
			err = UpdateNFT721Data(tokenID, contractAddr, from, to, chainID)
			if err != nil {
				logger.Logger.Error().Str("contractAddress", contractAddr).Int("chainID", chainID).
					Int("tokenID", tokenID).Str("from", from).Str("to", to).Msgf("Sync721EventData UpdateNFT721Data error:%s", err)
				if i == 0 {
					blockNum = 0
				} else {
					blockNum = logs[i-1].BlockNumber
				}
				err = errors.Wrap(err, "UpdateNFT721Data error")
				return
			}
			updateData := bson.D{
				{"$max", bson.D{
					{"BlockNum", int(log.BlockNumber)},
				}},
			}
			err := mongoData.MongoDB.UpdateOne(mongoData.DATABASE, mongoData.COLLECTION, filter, updateData)
			if err != nil {
				logger.Logger.Error().Str("contractAddress", contractAddr).Int("ChainID", chainID).Uint64("LatestBlockNum", log.BlockNumber).Msgf("Launch1155 error: %s", err)
			}
		}
	}
	logger.Logger.Info().Str("contractAddress", contractAddr).Int("chainID", chainID).
		Int64("from", fromBlock.Int64()).Int64("to", toBlock.Int64()).Msg("此721合约数据同步完毕")

	cache.SyncSignMap.Delete(contractAddr)
	if len(logs) == 0 {
		return toBlock.Uint64(), nil
	} else {
		return logs[len(logs)-1].BlockNumber, nil
	}
}

func commonArea(contractAddress common.Address, time int64, toBlock, fromBlock *big.Int, chainID int, eventSigns []common.Hash) []types.Log {
	fromBlockInt := fromBlock.Int64()
	logs := make([]types.Log, 0)
	indexAndLogsT := make(map[int64][]types.Log, 0)
	lock := sync.Mutex{}
	missedBlockNum := int64(-1)
	logger.Logger.Info().Str("contractAddress", util.StrToLow(contractAddress.String())).Int("chainID", chainID).
		Int64("start", fromBlock.Int64()).Int64("end", toBlock.Int64()).Msg("开始获取此合约此区间日志...")
	wg := sync.WaitGroup{}
	wg.Add(int(time))
	for i := int64(0); i < time; i++ {
		start := fromBlockInt + i*5000
		end := fromBlockInt + (i+1)*5000 - 1
		startBig := big.NewInt(start)
		endBig := big.NewInt(end)
		if i == time-1 {
			endBig = toBlock
		}
		go func(start, end *big.Int, i int64) {
			defer wg.Done()
			logsT, err := SyncContractEventLog(contractAddress, chainID, eventSigns, start, end)
			if err != nil {
				logger.Logger.Error().Str("contractAddress", util.StrToLow(contractAddress.String())).Int("chainID", chainID).
					Int64("start", start.Int64()).Int64("end", end.Int64()).Msgf("commonArea SyncContractEventLog error: %s", err)
				if missedBlockNum > start.Int64() || missedBlockNum == -1 {
					missedBlockNum = start.Int64()
				}
				return
			}
			logger.Logger.Info().Str("contractAddress", util.StrToLow(contractAddress.String())).Int("chainID", chainID).
				Int64("start", start.Int64()).Int64("end", end.Int64()).Int("此区间日志长度", len(logsT)).Msg("获取此合约此区间日志完毕")
			lock.Lock()
			indexAndLogsT[i] = logsT
			lock.Unlock()
		}(startBig, endBig, i)
	}
	wg.Wait()

	if missedBlockNum == -1 {
		for i := int64(0); i < time; i++ {
			logs = append(logs, indexAndLogsT[i]...)
		}
	} else {
		for i := int64(0); i < time; i++ {
			if len(indexAndLogsT[i]) > 0 && int64(indexAndLogsT[i][len(indexAndLogsT[i])-1].BlockNumber) < missedBlockNum {
				logs = append(logs, indexAndLogsT[i]...)
			}
		}
	}

	return logs
}

func SyncContractEventLog(contractAddr common.Address, chainID int, eventSigns []common.Hash, fromBlock, toBlock *big.Int) ([]types.Log, error) {
	//contractAddress := contractAddr.String()
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Topics:    [][]common.Hash{eventSigns},
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
		return logsOld, errors.Wrap(err, "client FilterLogs error")
	}
	return logsOld, nil
}

func DivideTime(fromBlock, toBlock *big.Int) int64 {
	fromBlockInt := fromBlock.Int64()
	toBlockInt := toBlock.Int64()
	var section = toBlockInt - fromBlockInt + 1
	var time = section / DIVID_CONTAIN
	if section%DIVID_CONTAIN > 0 {
		time++
	}
	return time
}
