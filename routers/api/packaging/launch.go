package packaging

import (
	"Ankr-gin-ERC721/mongoData"
	"Ankr-gin-ERC721/pkg/cache"
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	"Ankr-gin-ERC721/pkg/eth/subscribe"
	"Ankr-gin-ERC721/pkg/eth/syncData"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/msg"
	"Ankr-gin-ERC721/runtime"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/big"
	"net/http"
	"strconv"
)

func Launch721(contractAddr string, chainID int) (ok bool, httpCode, msgCode int, errorMsg string) {
	has := runtime.EventLoop.HasEvent(runtime.SUBSCRIBE_CONTRACT_721 + contractAddr)
	isSycning, ok := cache.SyncSignMap.Read(contractAddr)
	if !ok {
		isSycning = ok
	}
	if has || isSycning {
		return true, http.StatusOK, msg.SUCCESS, ""
	}
	cache.SyncSignMap.Insert(contractAddr)
	defer func() {
		if _,ok := cache.SyncSignMap[contractAddr];ok{
			cache.SyncSignMap.Delete(contractAddr)
		}
	}()
	filter := bson.D{{"ContractAddress", contractAddr}, {"ChainID", chainID}}
	latestBNum := cache.LatestBlockNum{}
	err := mongoData.MongoDB.FindOne(mongoData.DATABASE, mongoData.COLLECTION, filter, &latestBNum)
	if err != nil && err != mongo.ErrNoDocuments {
		return false, http.StatusInternalServerError, msg.ERROR_DB_ERC721, fmt.Sprintf("%s", errors.Wrap(err, "mongoDB FindOne error"))
	}
	mongoNull := err == mongo.ErrNoDocuments
	latestNum, err := LatestNum(chainID)
	if err != nil {
		return false, http.StatusInternalServerError, msg.ERROR_CLIENT_ERROR_ERC721, fmt.Sprintf("%s", errors.Wrap(err, "LatestNum error"))
	}
	var (
		start int64
		end   int64
	)

	if !mongoNull {
		start = int64(latestBNum.BlockNum + 1)
		txes, err := interactContract.GetNormalTXWithBlockNum(contractAddr, chainID, strconv.Itoa(int(start)), strconv.Itoa(int(latestNum)))
		if err == interactContract.NotFoundError {
			return true, http.StatusOK, msg.SUCCESS, ""
		}
		if err != nil {
			return false, http.StatusInternalServerError, msg.ERROR_ETHERSCAN_ERROR_ERC721, fmt.Sprintf("%s", errors.Wrap(err, "GetNormalTXWithBlockNum error"))
		}
		end_, _ := strconv.Atoi(txes[len(txes)-1].BlockNumber)
		end = int64(end_)
	} else {
		data := bson.M{"ContractAddress": contractAddr, "BlockNum": 0, "ChainID": chainID}
		err := mongoData.MongoDB.InsertOne(mongoData.DATABASE, mongoData.COLLECTION, data)
		if err != nil {
			logger.Logger.Error().Str("contractAddress", contractAddr).Int("ChainID", chainID).Int64("LatestBlockNum", end).Msgf("Launch721 error: %s", err)
		}

		txes, err := interactContract.GetAllNormalTX(contractAddr, chainID)
		if err == interactContract.NotFoundError {
			return true, http.StatusOK, msg.SUCCESS, ""
		}
		if err != nil {
			return false, http.StatusInternalServerError, msg.ERROR_ETHERSCAN_ERROR_ERC721, fmt.Sprintf("%s", errors.Wrap(err, "GetAllNormalTX error"))
		}
		start_, _ := strconv.Atoi(txes[0].BlockNumber)
		start = int64(start_)
		end = latestNum
	}

	blockNum, err := syncData.Sync721EventData(contractAddr, big.NewInt(start), big.NewInt(end), chainID)
	if err == nil {
		has = runtime.EventLoop.HasEvent(runtime.SUBSCRIBE_CONTRACT_721 + contractAddr)
		if !has {
			runtime.EventLoop.On(runtime.SUBSCRIBE_CONTRACT_721+contractAddr, subscribe.Subscribe721)
			go runtime.EventLoop.Emit(runtime.SUBSCRIBE_CONTRACT_721+contractAddr, subscribe.Data721{
				ContractAddr: contractAddr,
				ChainID:      chainID,
			})
		}
	}
	if err != nil {
		end = int64(blockNum)
		logger.Logger.Error().Str("contractAddress", contractAddr).Int("chainID", chainID).Msgf("Launch721 Sync721EventData error: %s", err)
	}

	return true, 0, 0, ""
}

func Launch1155(contractAddr string, chainID int) (ok bool, httpCode, msgCode int, errorMsg string) {
	has := runtime.EventLoop.HasEvent(runtime.SUBSCRIBE_CONTRACT_1155 + contractAddr)
	isSycning, ok := cache.SyncSignMap.Read(contractAddr)
	if !ok {
		isSycning = ok
	}
	if has || isSycning{
		return true, http.StatusOK, msg.SUCCESS, ""
	}
	cache.SyncSignMap.Insert(contractAddr)

	filter := bson.D{{"ContractAddress", contractAddr}, {"ChainID", chainID}}
	latestBNum := cache.LatestBlockNum{}
	err := mongoData.MongoDB.FindOne(mongoData.DATABASE, mongoData.COLLECTION, filter, &latestBNum)
	if err != nil && err != mongo.ErrNoDocuments {
		return false, http.StatusInternalServerError, msg.ERROR_DB_ERC721, fmt.Sprintf("%s", errors.Wrap(err, "mongoDB FindOne error"))
	}
	mongoNull := err == mongo.ErrNoDocuments
	latestNum, err := LatestNum(chainID)
	if err != nil {
		return false, http.StatusInternalServerError, msg.ERROR_CLIENT_ERROR_ERC721, fmt.Sprintf("%s", errors.Wrap(err, "LatestNum error"))
	}
	var (
		start int64
		end   int64
	)

	if !mongoNull {
		start = int64(latestBNum.BlockNum + 1)
		txes, err := interactContract.GetNormalTXWithBlockNum(contractAddr, chainID, strconv.Itoa(int(start)), strconv.Itoa(int(latestNum)))
		if err == interactContract.NotFoundError {
			return true, http.StatusOK, msg.SUCCESS, ""
		}
		if err != nil {
			return false, http.StatusInternalServerError, msg.ERROR_ETHERSCAN_ERROR_ERC721, fmt.Sprintf("%s", errors.Wrap(err, "GetNormalTXWithBlockNum error"))
		}
		end_, _ := strconv.Atoi(txes[len(txes)-1].BlockNumber)
		end = int64(end_)
	} else {
		data := bson.M{"ContractAddress": contractAddr, "BlockNum": 0, "ChainID": chainID}
		err := mongoData.MongoDB.InsertOne(mongoData.DATABASE, mongoData.COLLECTION, data)
		if err != nil {
			logger.Logger.Error().Str("contractAddress", contractAddr).Int("ChainID", chainID).Int64("LatestBlockNum", end).Msgf("Launch1155 error: %s", err)
		}

		txes, err := interactContract.GetAllNormalTX(contractAddr, chainID)
		if err == interactContract.NotFoundError {
			return true, http.StatusOK, msg.SUCCESS, ""
		}
		if err != nil {
			return false, http.StatusInternalServerError, msg.ERROR_ETHERSCAN_ERROR_ERC721, fmt.Sprintf("%s", errors.Wrap(err, "GetNormalTXWithBlockNum error"))
		}
		start_, _ := strconv.Atoi(txes[0].BlockNumber)
		start = int64(start_)
		end = latestNum
	}

	blockNum, err := syncData.Sync1155EventData(contractAddr, big.NewInt(start), big.NewInt(end), chainID)
	if err == nil {
		has = runtime.EventLoop.HasEvent(runtime.SUBSCRIBE_CONTRACT_1155 + contractAddr)
		if !has {
			runtime.EventLoop.On(runtime.SUBSCRIBE_CONTRACT_1155+contractAddr, subscribe.Subscribe1155)
			go runtime.EventLoop.Emit(runtime.SUBSCRIBE_CONTRACT_1155+contractAddr, subscribe.Data1155{
				ContractAddr: contractAddr,
				FromBlock:    big.NewInt(end),
				ChainID:      chainID,
			})
		}
	}
	if err != nil {
		end = int64(blockNum)
		logger.Logger.Error().Str("contractAddress", contractAddr).Int("chainID", chainID).Msgf("Launch1155 Sync1155EventData error: %s", err)
	}

	return true, 0, 0, ""
}
