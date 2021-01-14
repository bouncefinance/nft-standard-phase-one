package cache

import (
	"Ankr-gin-ERC721/mongoData"
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	"Ankr-gin-ERC721/pkg/msg"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func BaseURI1155Status(contractAddr string, chainID int)(ok bool,httpCode,msgCode int,bURI string) {
	bURI, err := interactContract.Get1155BaseURI(contractAddr, chainID)
	if err != nil {
		msg.MsgFlags[msg.ERROR_PARAM_ERC721] = fmt.Sprintf("There is an error with Get1155BaseURI. ERROR: %s", err)
		return false,http.StatusBadRequest,msg.ERROR_PARAM_ERC721,""
	}

	collection := mongoData.MongoDB.GetCollection(mongoData.DATABASE, mongoData.BASE_URI_1155_COLLECTION)
	filter := bson.D{{"ContractAddress", contractAddr}, {"ChainID", chainID}}
	baseURI := BaseURI1155{}
	err = collection.FindOne(context.Background(), filter).Decode(&baseURI)
	if err != nil && err != mongo.ErrNoDocuments {
		msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("There is an error with MongoDB. ERROR: %s", err)
		return false,http.StatusInternalServerError,msg.ERROR_DB_ERC721,""
	}
	if err == mongo.ErrNoDocuments {
		result, err := collection.InsertOne(context.Background(), bson.M{"ContractAddress": contractAddr, "BaseURI": bURI, "ChainID": chainID})
		if err != nil {
			fmt.Printf("GetMetadata1155 collection.InsertOne error: %v , result: %v \n", err, result)
			msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("There is an error with MongoDB. ERROR: %s", err)
			return false,http.StatusInternalServerError,msg.ERROR_DB_ERC721,""
		}
	} else if baseURI.BaseURI != bURI {
		updateData := bson.D{
			{"$set", bson.D{
				{"BaseURI", bURI},
			}},
		}
		updateResult, err := collection.UpdateOne(context.Background(), filter, updateData)
		if err != nil {
			fmt.Printf("collection.UpdateOne error: %v , result: %v \n", err, updateResult)
			msg.MsgFlags[msg.ERROR_DB_ERC721] = fmt.Sprintf("There is an error with MongoDB. ERROR: %s", err)
			return false,http.StatusInternalServerError,msg.ERROR_DB_ERC721,""
		}
	}
	return true,http.StatusOK,msg.SUCCESS,bURI
}
