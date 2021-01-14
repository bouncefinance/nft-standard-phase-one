package syncData

import (
	"Ankr-gin-ERC721/models"
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	logger "Ankr-gin-ERC721/pkg/logger"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

func UpdateSubscribe1155Data(tokenIDs []uint, amounts []*big.Int, contractAddr, from, to string, chainID int) error{
	tokens := make([]models.ERC1155, 0)
	newInt := big.NewInt(0)
	if from == common.BigToAddress(big.NewInt(0)).String() {
		for i := 0; i < len(tokenIDs); i++ {
			params := make(map[string]interface{})
			params["contract_addr"] = contractAddr
			params["token_id"] = tokenIDs[i]
			params["owner_addr"] = to
			erc1155sTo, err := models.GetERC1155(params)
			if err != nil {
				return errors.Wrap(err,"from == address(0) GetERC1155 error")
			}
			if len(erc1155sTo) > 0 {
				balance, _ := newInt.SetString(erc1155sTo[0].Balance, 0)
				balance.Add(balance, amounts[i])
				err := erc1155sTo[0].Update(map[string]interface{}{"balance": balance.String()})
				if err != nil {
					return errors.Wrap(err,"from == address(0) Update error")
				}
			} else {
				token := models.ERC1155{
					ContractAddr: contractAddr,
					TokenID:      tokenIDs[i],
					OwnerAddr:    to,
					Balance:      amounts[i].String(),
					ChainID:      chainID,
				}
				tokens = append(tokens, token)
			}
		}
		if len(tokens) > 0 {
			err:=models.AddERC1155(tokens)
			if err!=nil{
				return errors.Wrap(err,"from == address(0) AddERC1155 error")
			}
		}
	} else {
		for i := 0; i < len(tokenIDs); i++ {
			params := make(map[string]interface{})
			params["contract_addr"] = contractAddr
			params["token_id"] = tokenIDs[i]
			params["owner_addr"] = from
			erc1155sFrom, err := models.GetERC1155(params)
			if err != nil {
				return errors.Wrap(err,"from != address(0) GetERC1155 from error")
			}

			params["owner_addr"] = to
			erc1155sTo, err := models.GetERC1155(params)
			if err != nil {
				return errors.Wrap(err,"from != address(0) GetERC1155 to error")
			}
			if len(erc1155sTo) == 0 {
				token := models.ERC1155{
					ContractAddr: contractAddr,
					TokenID:      tokenIDs[i],
					OwnerAddr:    to,
					Balance:      amounts[i].String(),
					ChainID:      chainID,
				}
				tokens = append(tokens, token)
			} else {
				balance, _ := newInt.SetString(erc1155sTo[0].Balance, 0)
				balance.Add(balance, amounts[i])
				err := erc1155sTo[0].Update(map[string]interface{}{"balance": balance.String()})
				if err != nil {
					return errors.Wrap(err,"from != address(0) Update to error")
				}
			}
			if len(erc1155sFrom) == 0 {
				//logger.Logger.Info().Str("contractAddress",contractAddr).Str("from",from).Str("to",to).Uint("tokenID",tokenIDs[i]).Msg("from does not have this ID tokens.")
				return errors.Wrap(fmt.Errorf("from is not address(0) but has not this 1155 ID token"),"")
			}
			balance, _ := newInt.SetString(erc1155sFrom[0].Balance, 0)
			balance.Sub(balance, amounts[i])
			err = erc1155sFrom[0].Update(map[string]interface{}{"balance": balance.String()})
			if err != nil {
				return errors.Wrap(err,"from != address(0) Update from error")
			}
		}
		if len(tokens) > 0 {
			err:=models.AddERC1155(tokens)
			if err!=nil{
				return errors.Wrap(err,"from != address(0) AddERC1155 error")
			}
		}
	}
	return nil
}

func Update1155URI(contractAddr string, tokenID int, uri string)error {
	erc1155URI, err := models.GetERC1155URI(map[string]interface{}{"contract_addr": contractAddr, "token_id": tokenID})
	if err != nil {
		return errors.Wrap(err,"GetERC1155URI error")
	}
	if len(erc1155URI) > 0 {
		err := erc1155URI[0].Update(map[string]interface{}{"token_uri": erc1155URI[0].TokenURI})
		if err != nil {
			return errors.Wrap(err,"Update error")
		}
	} else {
		uri1155 := &models.ERC1155URI{
			ContractAddr: contractAddr,
			TokenID:      tokenID,
			TokenURI:     uri,
		}
		err := uri1155.Add()
		if err != nil {
			return errors.Wrap(err,"Add error")
		}
	}
	return nil
}

func UpdateNFT721Data(tokenID int, contractAddr, from, to string, chainID int) error{
	params := make(map[string]interface{})
	params["contract_addr"] = contractAddr
	params["token_id"] = tokenID
	params["owner_addr"] = from

	token := &models.ERC721{}
	if from != common.BigToAddress(big.NewInt(0)).String() {
		newData := make(map[string]interface{})
		condition := make(map[string]interface{})
		condition["contract_addr"] = params["contract_addr"]
		condition["token_id"] = params["token_id"]
		newData["owner_addr"] = to
		err := token.UpdateByCondition(condition, newData)
		if err != nil {
			return errors.Wrap(err,"UpdateByCondition error")
		}
		logger.Logger.Info().Str("contractAddress",contractAddr).Int("chainID",chainID).
			Int("tokenID",tokenID).Str("from",from).Str("to",to).Msg("721 nft data has updated successfully")
	} else {
		uri, err := interactContract.GetTokenURIStr(contractAddr, chainID, tokenID)
		if err != nil {
			logger.Logger.Error().Str("contractAddress",contractAddr).Int("chainID",chainID).
				Int("tokenID",tokenID).Msgf("UpdateNFT721Data GetTokenURIStr error:%s",err)
		}
		tokens := make([]models.ERC721, 0)
		tokens = append(tokens, models.ERC721{
			ContractAddr: contractAddr,
			TokenID:      tokenID,
			OwnerAddr:    to,
			ChainID:      chainID,
			TokenURI:     uri,
		})
		err = models.AddERC721(tokens)
		if err != nil {
			return errors.Wrap(err,"AddERC721 error")
		}
		logger.Logger.Info().Str("contractAddress",contractAddr).Int("chainID",chainID).
			Int("tokenID",tokenID).Str("to",to).Msg("721 nft data has added successfully")
	}
	return nil
}
