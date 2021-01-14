package api

import (
	"Ankr-gin-ERC721/models"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/msg"
	"Ankr-gin-ERC721/pkg/util"
	"Ankr-gin-ERC721/routers/api/packaging"
	"Ankr-gin-ERC721/runtime"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetERC1155(c *gin.Context) {
	appC := runtime.Context{C: c}
	valid := &validation.Validation{}
	userAddr_, contractAddr_, chainID := validateParam(c, ERC1155_URL, valid)
	if valid.HasErrors() {
		msg.MsgFlags[msg.ERROR_PARAM_ERC721] = fmt.Sprintf("%s", valid.Errors)
		appC.Response(http.StatusBadRequest, msg.ERROR_PARAM_ERC721, nil)
		return
	}

	userAddr := util.StrToLow(userAddr_)
	contractAddr := util.StrToLow(contractAddr_)
	tokens := make([]models.ERC1155, 0)

	logger.Logger.Info().Str("userAddr",userAddr).Str("contractAddr",contractAddr).Int("chainID",chainID).Msg("user request GetERC1155.")

	/*start, end, err := packaging.StartAndEndNum(contractAddr, chainID)
	if err != nil && err != interactContract.NotFoundError {
		logger.Logger.Error().Str("userAddress", userAddr).Str("contractAddress", contractAddr).Int("chainID",chainID).Msgf("GetERC1155 StartAndEndNum error: %s",err)
		msg.MsgFlags[msg.ERROR_RPC_ERROR_ERC721] = fmt.Sprintf("GetERC721 StartAndEndNum error: %s", err)
		appC.Response(http.StatusInternalServerError, msg.ERROR_RPC_ERROR_ERC721, nil)
		return
	}
	if err == interactContract.NotFoundError {
		appC.Response(http.StatusOK, msg.SUCCESS, nil)
		return
	}*/

	ok, code, msgCode ,errMsg:= packaging.Launch1155(contractAddr, chainID)
	if !ok{
		logger.Logger.Error().Str("userAddr",userAddr).Str("contractAddr", contractAddr).Int("chainID",chainID).Msgf("GetERC1155 Launch1155 error:%s", errMsg)
		msg.MsgFlags[msgCode] = errMsg
		appC.Response(code, msgCode, nil)
		return
	}

	tokens, err := models.GetERC1155(map[string]interface{}{"contract_addr": contractAddr, "owner_addr": userAddr})
	if err != nil {
		logger.Logger.Error().Str("userAddr",userAddr).Str("contractAddr", contractAddr).Int("chainID",chainID).Msgf("GetERC1155 models.GetERC1155 error: %s",err)
		msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("models.GetERC721 error: %s", err)
		appC.Response(http.StatusBadRequest, msg.ERROR_DB_ERC721, nil)
		return
	}

	data := make(map[string]interface{})
	data["tokens"] = tokens
	appC.Response(http.StatusOK, msg.SUCCESS, data)
}

