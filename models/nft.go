package models

type NFT struct {
	ContractAddr string `json:"contract_addr" 	gorm:"column:contract_addr"`
	TokenType    string `json:"token_type" 		gorm:"column:token_type"`
	TokenID      int    `json:"token_id" 		gorm:"column:token_id"`
	OwnerAddr    string `json:"owner_addr" 		gorm:"column:owner_addr"`
	ChainID      int    `json:"chain_id" 		gorm:"column:chain_id"`
	Balance      string `json:"balance" 		gorm:"column:balance"`
	TokenURI     string `json:"token_uri" 		gorm:"column:token_uri"`

	Name        interface{} `json:"name"`
	Description interface{} `json:"description"`
	Image       interface{} `json:"image"`
}

