package models

import (
	"Ankr-gin-ERC721/pkg/logger"
	"fmt"
	"gorm.io/gorm"
)

type ERC721 struct {
	gorm.Model
	ContractAddr string `json:"contract_addr" 	gorm:"column:contract_addr"`
	TokenID      int    `json:"token_id" 		gorm:"column:token_id"`
	OwnerAddr    string `json:"owner_addr" 		gorm:"column:owner_addr"`
	ChainID      int    `json:"chain_id" 		gorm:"column:chain_id"`
	TokenURI     string `json:"token_uri" 		gorm:"column:token_uri"`
}

func (erc721 *ERC721) Update(newData map[string]interface{}) error {
	return db.Model(erc721).Updates(newData).Error
}
func (erc721 *ERC721) UpdateByCondition(condition, newData map[string]interface{}) error {
	return db.Model(erc721).Where(condition).Updates(newData).Error
}

func GetERC721(maps map[string]interface{}) (tokens []ERC721, err error) {
	err = db.Where(maps).Find(&tokens).Error
	return
}

func AddERC721(erc721s []ERC721) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover error: ", r)
			tx.Rollback()
		}
	}()

	tx.Create(&erc721s)

	for _, erc721 := range erc721s {
		logger.Logger.Info().Str("contract address",erc721.ContractAddr).Int("chainID",erc721.ChainID).
			Int("erc721 len", len(erc721s)).Str("owner",erc721.OwnerAddr).Int("tokenID",erc721.TokenID).Msg("数据库事务提交，erc721")
	}
	return tx.Commit().Error
}
