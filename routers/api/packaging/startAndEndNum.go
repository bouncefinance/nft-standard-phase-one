package packaging

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	"Ankr-gin-ERC721/pkg/setting"
	"context"
	"github.com/pkg/errors"
	"strconv"
)

func StartAndEndNum(address string, chainID int) (int64, int64, error) {
	time := 0
	txes, err := interactContract.GetAllNormalTX(address, chainID)
	for err != nil && time < conf.RETRY_TIME {
		txes, err = interactContract.GetAllNormalTX(address, chainID)
		time++
	}
	if err != nil && err != interactContract.NotFoundError {
		return -1, -1, errors.Wrap(err, "GetAllNormalTX error")
	}
	if len(txes) == 0 {
		return -1, -1, interactContract.NotFoundError
	}
	startNumStr := txes[0].BlockNumber
	startNum, _ := strconv.ParseInt(startNumStr, 0, 0)

	time = 0
	endBlockHeader, err := setting.ETHHTTPClients[chainID].HeaderByNumber(context.Background(), nil)
	for err != nil && time < conf.RETRY_TIME {
		endBlockHeader, err = setting.ETHHTTPClients[chainID].HeaderByNumber(context.Background(), nil)
		time++
	}
	if err != nil {
		return -1, -1, errors.Wrap(err, "client.HeaderByNumber error")
	}
	endNum := endBlockHeader.Number.Int64()
	return startNum, endNum, nil
}

func LatestNum(chainID int) (int64, error) {
	time := 0
	endBlockHeader, err := setting.ETHHTTPClients[chainID].HeaderByNumber(context.Background(), nil)
	for err != nil && time < conf.RETRY_TIME {
		endBlockHeader, err = setting.ETHHTTPClients[chainID].HeaderByNumber(context.Background(), nil)
		time++
	}
	if err != nil {
		return -1,errors.Wrap(err, "LatestNum error")
	}
	return endBlockHeader.Number.Int64(),err
}