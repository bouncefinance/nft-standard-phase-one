package api

import (
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/setting"

	"context"
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

const (
	ERC721_URL  = "erc721"
	ERC1155_URL = "erc1155"
)

func validateParam(c *gin.Context, url string, valid *validation.Validation) (userAddr string, contractAddr string, chainID int) {
	userAddr = c.Query("user_addr")
	contractAddr = c.Query("contract_addr")
	chainID = com.StrTo(c.Query("chain_id")).MustInt()

	valid.Required(userAddr, "user_addr").Message("error: userAddr is null!")
	valid.Required(contractAddr, "contract_addr").Message("error: contractAddr is null!")
	valid.Required(chainID, "chain_id").Message("error: chainID is null!")

	if len(userAddr) != addr_len || len(contractAddr) != addr_len {
		valid.Error("the length of user's or contract's address is invalid.")
		return
	}
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(userAddr) || !re.MatchString(contractAddr) {
		valid.Error("please enter correct address.")
		return
	}

	_, ok := setting.ETHClients[chainID]
	if !ok {
		valid.Error("chainID is not correct.")
		return
	}

	address := common.HexToAddress(contractAddr)
	setting.ClientLock.Lock()
	byteCode, err := setting.ETHClients[chainID].CodeAt(context.Background(), address, nil) // nil is latest block
	setting.ClientLock.Unlock()
	if err != nil {
		valid.Error("ETHClient.CodeAt error: %s", err)
		return
	}
	if !(len(byteCode) > 0) {
		valid.Error("Input content is not contractAddress.")
		return
	}

	var (
		interSign string
	)
	if url == ERC721_URL {
		interSign = interactContract.ERC721_INTERFACE_ID
	} else if url == ERC1155_URL {
		interSign = interactContract.ERC1155_INTERFACE_ID
	}
	errV := IsSupportInterface(contractAddr, interSign, chainID)
	if errV != nil {
		valid.Errors = append(valid.Errors, errV)
	}
	return
}

func IsSupportInterface(contractAddr string, interSign string, chainID int) *validation.Error {
	data, err := interactContract.SupportInterface(contractAddr, interSign, chainID)
	if err != nil {
		logger.Logger.Error().Str("contractAddress",contractAddr).Int("chainID",chainID).Msgf("IsSupportInterface SupportInterface error: %s",err)
		return &validation.Error{
			Message: "There is a error with ETH node OR Contract is not correct contractAddress.",
		}
	}
	returnData := interactContract.ReturnData{}
	err = json.Unmarshal(data, &returnData)
	if err != nil {
		logger.Logger.Error().Str("contractAddress",contractAddr).Int("chainID",chainID).Msgf("IsSupportInterface json.Unmarshal error: %s",err)
		return &validation.Error{
			Message: "Contract is not correct contractAddress.",
		}
	}
	i, err := strconv.ParseInt(returnData.Result, 0, 0)
	if err != nil {
		logger.Logger.Error().Str("contractAddress",contractAddr).Int("chainID",chainID).Msgf("IsSupportInterface strconv.ParseInt error: %s",err)
		return &validation.Error{
			Message: "Contract is not correct contractAddress.",
		}
	}
	if i != interactContract.ACTIVELY_RESPONSE {
		return &validation.Error{
			Message: "Contract is not correct contractAddress.",
		}
	}
	return nil
}
