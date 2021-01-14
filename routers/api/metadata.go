package api

import (
	"Ankr-gin-ERC721/models"
	"Ankr-gin-ERC721/pkg/cache"
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/msg"
	"Ankr-gin-ERC721/pkg/util"
	"Ankr-gin-ERC721/routers/api/packaging"
	"Ankr-gin-ERC721/runtime"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

const ()

func GetMetadata721(c *gin.Context) {
	appC := runtime.Context{C: c}
	contractAddr := c.Query("contract_addr")
	chainID := com.StrTo(c.Query("chain_id")).MustInt()
	limit := com.StrTo(c.Query("limit")).MustInt()
	valid := &validation.Validation{}
	valid.Required(contractAddr, "contract_addr").Message("error: contractAddr is null!")
	valid.Required(chainID, "chain_id").Message("error: chainID is null!")
	//valid.Required(limit, "limit").Message("error: limit is null!")

	//valid.Min(limit, 0, "limit").Message("error: limit is too small!")

	vErr := IsSupportInterface(contractAddr, interactContract.ERC721_INTERFACE_ID, chainID)
	if vErr != nil {
		valid.Errors = append(valid.Errors, vErr)
	}
	if valid.HasErrors() {
		msg.MsgFlags[msg.ERROR_PARAM_ERC721] = fmt.Sprintf("%s", valid.Errors)
		appC.Response(http.StatusBadRequest, msg.ERROR_PARAM_ERC721, nil)
		return
	}

	contractAddr = util.StrToLow(contractAddr)
	logger.Logger.Info().Int("limit", limit).Str("contractAddr", contractAddr).Int("chainID", chainID).Msg("user request GetMetadata721.")

	/*txes, err := interactContract.GetAllNormalTX(contractAddr, chainID)
	if err != nil && err != interactContract.NotFoundError{
		msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("runtime.GetAllTXFromContract error: %s", err)
		appC.Response(http.StatusBadRequest, msg.ERROR_DB_ERC721, nil)
		return
	}
	startNumStr := txes[0].BlockNumber
	endNumStr := txes[len(txes)-1].BlockNumber
	startNum, _ := strconv.ParseInt(startNumStr, 0, 0)
	endNum, _ := strconv.ParseInt(endNumStr, 0, 0)*/

	ok, code, msgCode, errMsg := packaging.Launch721(contractAddr, chainID)
	if !ok {
		logger.Logger.Error().Int("limit", limit).Str("contractAddr", contractAddr).Int("chainID", chainID).Msgf("GetMetadata721 Launch721 error:%s", errMsg)
		msg.MsgFlags[msgCode] = errMsg
		appC.Response(code, msgCode, nil)
		return
	}

	tokens, err := models.GetERC721(map[string]interface{}{"contract_addr": contractAddr, "chain_id": chainID})
	if err != nil {
		logger.Logger.Error().Int("limit", limit).Str("contractAddr", contractAddr).Int("chainID", chainID).Msgf("GetMetadata721 models.GetERC721 error: %s", err)
		msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("models.GetERC721 error: %s", err)
		appC.Response(http.StatusBadRequest, msg.ERROR_DB_ERC721, nil)
		return
	}
	if 0 < limit && limit < len(tokens) {
		tokens = tokens[:limit]
	}

	result := make(map[int]interface{}, 0)
	for _, token := range tokens {
		var dataI interface{}
		data, err := util.GetUrl(token.TokenURI, nil)
		if err != nil {
			logger.Logger.Error().Int("limit", limit).Str("contractAddr", contractAddr).Int("chainID", chainID).Msgf("GetMetadata721 util.GetUrl error: %s", err)
			continue
		}
		err = json.Unmarshal(data, &dataI)
		if err != nil {
			logger.Logger.Error().Int("limit", limit).Str("contractAddr", contractAddr).Int("chainID", chainID).Msgf("GetMetadata721 json.Unmarshal error: %s", err)
			continue
		}
		d, ok := dataI.(map[string]interface{})
		if !ok {
			continue
		}
		d1, ok := d["properties"].(map[string]interface{})
		if !ok {
			continue
		}
		dName, ok := d1["name"].(map[string]interface{})
		if !ok {
			continue
		}
		dDescription, ok := d1["description"].(map[string]interface{})
		if !ok {
			continue
		}
		dImage, ok := d1["image"].(map[string]interface{})
		if !ok {
			continue
		}
		r := make(map[string]interface{})
		r["name"] = dName["description"]
		r["description"] = dDescription["description"]
		r["image"] = dImage["description"]

		result[token.TokenID] = r
	}

	appC.ResponseMetaData(http.StatusOK, msg.SUCCESS, result)
}

func GetMetadata1155(c *gin.Context) {
	appC := runtime.Context{C: c}
	contractAddr := c.Query("contract_addr")
	chainID := com.StrTo(c.Query("chain_id")).MustInt()
	limit := com.StrTo(c.Query("limit")).MustInt()

	valid := &validation.Validation{}
	valid.Required(contractAddr, "contract_addr").Message("error: contractAddr is null!")
	valid.Required(chainID, "chain_id").Message("error: chainID is null!")

	vErr := IsSupportInterface(contractAddr, interactContract.ERC1155_INTERFACE_ID, chainID)
	if vErr != nil {
		valid.Errors = append(valid.Errors, vErr)
	}
	if valid.HasErrors() {
		msg.MsgFlags[msg.ERROR_PARAM_ERC721] = fmt.Sprintf("%s", valid.Errors)
		appC.Response(http.StatusBadRequest, msg.ERROR_PARAM_ERC721, nil)
		return
	}
	contractAddr = util.StrToLow(contractAddr)

	logger.Logger.Info().Int("limit", limit).Str("contractAddr", contractAddr).Int("chainID", chainID).Msg("user request GetMetadata1155.")
	ok, httpCode, msgCode, bURI := cache.BaseURI1155Status(contractAddr, chainID)
	if !ok {
		appC.Response(httpCode, msgCode, nil)
		return
	}

	ok, code, msgCode, errMsg := packaging.Launch1155(contractAddr, chainID)
	if !ok {
		logger.Logger.Error().Int("limit", limit).Str("contractAddr", contractAddr).Int("chainID", chainID).Msgf("GetMetadata1155 Launch1155 error:%s", errMsg)
		msg.MsgFlags[msgCode] = errMsg
		appC.Response(code, msgCode, nil)
		return
	}

	tokens1155, err := models.GetERC1155URI(map[string]interface{}{"contract_addr": contractAddr})
	if err != nil {
		logger.Logger.Error().Int("limit", limit).Str("contractAddr", contractAddr).Int("chainID", chainID).Msgf("GetMetadata1155 models.GetERC1155URI error: %s", err)
		msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("models.GetERC1155URI error: %s", err)
		appC.Response(http.StatusBadRequest, msg.ERROR_DB_ERC721, nil)
		return
	}
	result := make(map[int]interface{}, 0)
	if len(tokens1155) > 0 {
		for i := 0; ; {
			if i >= len(tokens1155)-1 {
				break
			}
			if tokens1155[i+1].TokenID == tokens1155[i].TokenID {
				tokens1155 = append(tokens1155[:i], tokens1155[i+1:]...)
				continue
			}
			i++
		}
		if 0 < limit && limit < len(tokens1155) {
			tokens1155 = tokens1155[:limit]
		}

		for _, token := range tokens1155 {
			var dataI interface{}
			data,err := util.GetUrl(bURI+token.TokenURI, nil)
			if err != nil {
				logger.Logger.Error().Int("tokenID", token.TokenID).Str("contractAddr", contractAddr).Int("chainID", chainID).Msgf("GetMetadata1155 util.GetUrl error: %s", err)
				continue
			}
			err = json.Unmarshal(data, &dataI)
			if err != nil {
				logger.Logger.Error().Int("tokenID", token.TokenID).Str("contractAddr", contractAddr).Int("chainID", chainID).Msgf("GetMetadata1155 json.Unmarshal error: %s", err)
				continue
			}
			result[token.TokenID] = dataI
		}
	}

	appC.ResponseMetaData(http.StatusOK, msg.SUCCESS, result)
}
