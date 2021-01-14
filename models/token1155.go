package models

import (
	"fmt"
	"gorm.io/gorm"
)

type ERC1155 struct {
	gorm.Model
	ContractAddr string `json:"contract_addr" 	gorm:"column:contract_addr"`
	TokenID      uint    `json:"token_id" 		gorm:"column:token_id"`
	OwnerAddr    string `json:"owner_addr" 		gorm:"column:owner_addr"`
	Balance      string `json:"balance" 		gorm:"column:balance"`
	ChainID      int    `json:"chain_id" 		gorm:"column:chain_id"`
	//TokenURI     string `json:"token_uri" 		gorm:"column:token_uri"`
}

func (erc1155 *ERC1155) Update(newData map[string]interface{}) error {
	return db.Model(erc1155).Updates(newData).Error
}

func GetERC1155(maps map[string]interface{}) (tokens []ERC1155, err error) {
	err = db.Where(maps).Find(&tokens).Error
	return
}

func AddERC1155(erc1155s []ERC1155) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover error: ", r)
			tx.Rollback()
		}
	}()

	tx.Create(&erc1155s)

	return tx.Commit().Error
}

type ERC1155URI struct {
	gorm.Model
	ContractAddr string `json:"contract_addr" 	gorm:"column:contract_addr"`
	TokenID      int    `json:"token_id" 		gorm:"column:token_id"`
	TokenURI     string `json:"token_uri" 		gorm:"column:token_uri"`
}

func GetERC1155URI(maps map[string]interface{}) (uris []ERC1155URI, err error) {
	err = db.Where(maps).Find(&uris).Error
	return
}

func (uri *ERC1155URI) Update(newData map[string]interface{}) error {
	return db.Model(uri).Updates(newData).Error
}
func (uri *ERC1155URI) Add() error {
	return db.Create(uri).Error
}
