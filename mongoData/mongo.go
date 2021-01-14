package mongoData

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/logger"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoDatabase struct {
	MongoClient *mongo.Client
}

func (m *MongoDatabase) GetCollection(database string, collection string) *mongo.Collection {
	return m.MongoClient.Database(database).Collection(collection)
}

func (m *MongoDatabase) FindOne(database string, collection_ string, filter bson.D, object interface{})error {
	time:=0
	collection := m.GetCollection(database, collection_)
	err := collection.FindOne(context.Background(), filter).Decode(object)
	for err != nil && err != mongo.ErrNoDocuments && time < conf.RETRY_TIME {
		err = collection.FindOne(context.Background(), filter).Decode(object)
		time++
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err,"mongo FindOne error")
	}else if err == mongo.ErrNoDocuments{
		return err
	}
	return nil
}

func (m *MongoDatabase) InsertOne(database string, collection_ string, data bson.M)error {
	time := 0
	collection := m.GetCollection(database, collection_)
	_, err := collection.InsertOne(context.Background(), data)
	for err != nil && time < conf.RETRY_TIME {
		_, err = collection.InsertOne(context.Background(), data)
		time++
	}
	if err != nil {
		return errors.Wrap(err,"mongo InsertOne error")
	}
	return nil
}

func (m *MongoDatabase) UpdateOne(database string, collection_ string,filter bson.D, data bson.D)error {
	time := 0
	collection := m.GetCollection(database, collection_)
	_, err := collection.UpdateOne(context.Background(), filter, data)
	for err != nil && time < conf.RETRY_TIME {
		_, err = collection.UpdateOne(context.Background(), filter, data)
		time++
	}
	if err != nil {
		return errors.Wrap(err,"mongo UpdateOne error")
	}
	return nil
}


var MongoDB *MongoDatabase

const (
	DATABASE                 = "AnkrNFT"
	COLLECTION               = "LatestBlockNum"
	BASE_URI_1155_COLLECTION = "BaseURI1155"
	ADDRESS_NFT_COLLECTION   = "AddressNFT"
)

func init() {
	MongoDB = &MongoDatabase{
		MongoClient: SetConnect(),
	}
}

// 连接设置
func SetConnect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MONGO_UTL).SetMaxPoolSize(20)) // 连接池
	if err != nil {
		log.Fatal("mongo connect failed error: ", err)
	}
	logger.Logger.Info().Str("mongodb link", conf.MONGO_UTL).Msg("")
	return client
}
