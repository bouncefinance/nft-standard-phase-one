package api

import (
	"Ankr-gin-ERC721/models"
	"Ankr-gin-ERC721/mongoData"
	"Ankr-gin-ERC721/pkg/cache"
	"Ankr-gin-ERC721/pkg/logger"
	"Ankr-gin-ERC721/pkg/msg"
	nfts2 "Ankr-gin-ERC721/pkg/nfts"
	"Ankr-gin-ERC721/pkg/setting"
	"Ankr-gin-ERC721/pkg/util"
	"Ankr-gin-ERC721/runtime"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/unknwon/com"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"regexp"
)

func GetAllNFT(c *gin.Context) {
	appC := runtime.Context{C: c}
	valid := &validation.Validation{}

	address := c.Query("address")
	chainID := com.StrTo(c.Query("chain_id")).MustInt()
	address = util.StrToLow(address)
	valid.Required(address, "address").Message("error: address is null!")
	valid.Required(chainID, "chain_id").Message("error: chainID is null!")

	if len(address) != addr_len {
		valid.Error("the length of user's address is invalid.")
		return
	}
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(address) {
		valid.Error("please enter correct address.")
		return
	}
	//	判断chainID是否被程序支持
	_, ok := setting.ETHClients[chainID]
	if !ok {
		valid.Error("chainID is not correct.")
		return
	}
	if valid.HasErrors() {
		msg.MsgFlags[msg.ERROR_PARAM_ERC721] = fmt.Sprintf("%s", valid.Errors)
		appC.Response(http.StatusBadRequest, msg.ERROR_PARAM_ERC721, nil)
		return
	}

	logger.Logger.Info().Str("userAddress", address).Int("chainID", chainID).Msg("user request GetAllNFT.")

	filter := bson.D{{"Address", address}, {"ChainID", chainID}}
	baseURI := cache.BaseURI1155{}
	err := mongoData.MongoDB.FindOne(mongoData.DATABASE, mongoData.ADDRESS_NFT_COLLECTION, filter, &baseURI)
	if err != nil && err != mongo.ErrNoDocuments {
		logger.Logger.Error().Str("userAddress", address).Int("chainID", chainID).Msgf("GetAllNFT mongo FindOne error: %s", err)
		msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("There is an error with MongoDB. ERROR: %s", err)
		appC.Response(http.StatusInternalServerError, msg.ERROR_DB_ERC721, nil)
		return
	}

	err = nfts2.All721NFTAssets(address, chainID)
	if err != nil {
		logger.Logger.Error().Str("userAddress", address).Int("chainID", chainID).Msgf("GetAllNFT All721NFTAssets error: %s", err)
		msg.MsgFlags[msg.ERROR_ETHERSCAN_ERROR_ERC721] = fmt.Sprintf("There is an error with Etherscan api. ERROR: %s", err)
		appC.Response(http.StatusInternalServerError, msg.ERROR_ETHERSCAN_ERROR_ERC721, nil)
		return
	}

	if err == mongo.ErrNoDocuments {
		err = mongoData.MongoDB.InsertOne(mongoData.DATABASE, mongoData.ADDRESS_NFT_COLLECTION, bson.M{"Address": address, "ChainID": chainID})
		if err != nil {
			logger.Logger.Error().Str("userAddress", address).Int("chainID", chainID).Msgf("GetAllNFT mongo InsertOne error: %s", err)
		}
	}

	nfts := make([]models.NFT, 0)
	tokens721, err := models.GetERC721(map[string]interface{}{"owner_addr": address, "chain_id": chainID})
	if err != nil {
		logger.Logger.Error().Str("userAddress", address).Int("chainID", chainID).Msgf("GetAllNFT models.GetERC721 error: %s", err)
		msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("models.GetERC721 error: %s", err)
		appC.Response(http.StatusBadRequest, msg.ERROR_DB_ERC721, nil)
		return
	}
	for _, erc721 := range tokens721 {
		n := models.NFT{
			ContractAddr: erc721.ContractAddr,
			TokenType:    "721",
			TokenID:      erc721.TokenID,
			OwnerAddr:    erc721.OwnerAddr,
			ChainID:      erc721.ChainID,
			Balance:      "1",
			TokenURI:     erc721.TokenURI,
		}

		r, err := GetMeta721(&erc721)
		if err != nil {
			logger.Logger.Error().Int("token id", erc721.TokenID).Str("contract address", erc721.ContractAddr).
				Int("chainID", chainID).Msgf("GetAllNFT GetMeta721 error: %s", err)
		}

		n.Name = r["name"]
		n.Description = r["description"]
		n.Image = r["image"]
		nfts = append(nfts, n)
	}
	/*tokens1155, err := models.GetERC1155(map[string]interface{}{"owner_addr": address, "chain_id": chainID})
	if err != nil {
		msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("models.GetERC721 error: %s", err)
		appC.Response(http.StatusBadRequest, msg.ERROR_DB_ERC721, nil)
		return
	}
	for _, erc1155 := range tokens1155 {
		//	获取base URI
		ok, _, _, bURI := cache.BaseURI1155Status(erc1155.ContractAddr, chainID)
		if !ok {
			fmt.Printf("cache.BaseURI1155Status err ContractAddr: %s chainID: %d\n",erc1155.ContractAddr,chainID)
			continue
		}
		//	获取token uri
		uriData, err := models.GetERC1155URI(map[string]interface{}{"contract_addr": erc1155.ContractAddr, "token_id": erc1155.TokenID})
		if err != nil {
			fmt.Println("models.GetERC1155URI error:", err)
			continue
		}
		nfts = append(nfts, models.NFT{
			ContractAddr: erc1155.ContractAddr,
			TokenType:    "1155",
			TokenID:      erc1155.TokenID,
			OwnerAddr:    erc1155.OwnerAddr,
			ChainID:      chainID,
			Balance:      erc1155.Balance,
			TokenURI:     bURI + uriData[0].TokenURI,
		})
	}*/

	data := make(map[string]interface{})
	data["nfts"] = nfts
	appC.Response(http.StatusOK, msg.SUCCESS, data)
}

func GetMeta721(token *models.ERC721) (map[string]interface{}, error) {
	var dataI interface{}
	data, err := util.GetUrl(token.TokenURI, nil)
	if err != nil {
		return nil, errors.Wrap(err, "util.GetUrl error")
	}
	err = json.Unmarshal(data, &dataI)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal error")
	}
	d, ok := dataI.(map[string]interface{})
	if !ok {
		return nil, errors.New("dataI type error")
	}
	d1, ok := d["properties"].(map[string]interface{})
	if !ok {
		return nil, errors.New("d[\"properties\"] type error")
	}
	dName, ok := d1["name"].(map[string]interface{})
	if !ok {
		return nil, errors.New("d1[\"name\"] type error")
	}
	dDescription, ok := d1["description"].(map[string]interface{})
	if !ok {
		return nil, errors.New("d1[\"description\"] type error")
	}
	dImage, ok := d1["image"].(map[string]interface{})
	if !ok {
		return nil, errors.New("d1[\"image\"] type error")
	}
	r := make(map[string]interface{})
	r["name"] = dName["description"]
	r["description"] = dDescription["description"]
	r["image"] = dImage["description"]
	return r, nil
}
