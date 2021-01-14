package models

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
)

func TestAddERC721(t *testing.T) {
	tokens := []ERC721{
		{
			ContractAddr: "0x93e508f373690cC4307a7A2363e573E63dAEF54E",
			TokenID:      3,
			OwnerAddr:    "0xbccc2073adfc46421308f62cfd9868df00d339a8",
			TokenURI:     "www.baidu.com",
		},
		{
			ContractAddr: "0x93e508f373690cC4307a7A2363e573E63dAEF54E",
			TokenID:      4,
			OwnerAddr:    "0xbccc2073adfc46421308f62cfd9868df00d339a8",
		},
	}
	err := AddERC721(tokens)
	if err != nil {
		t.Error("error: ", err)
	}
}

func TestERC721_Update(t *testing.T) {
	var token = &ERC721{
		Model:        gorm.Model{ID: 21},
		ContractAddr: "0x93e508f373690cC4307a7A2363e573E63dAEF54E",
		TokenID:      17,
		OwnerAddr:    "0x8eedc16d4921c827451dd",
	}

	newData := make(map[string]interface{})
	newData["token_id"] = 16
	err := token.Update(newData)
	if err != nil {
		t.Error("token.Update error ==>", err)
	}
}

func TestERC721_UpdateByCondition(t *testing.T) {
	var token = &ERC721{}

	newData := make(map[string]interface{})
	condition := make(map[string]interface{})

	condition["contract_addr"] = "0xd8d638be21b4101e1858de84d4540aed2d02674d"
	condition["token_id"] = 2
	newData["owner_addr"] = "0xBCcC2073ADfC46421308f62cfD9868dF00D339a8"

	err := token.UpdateByCondition(condition, newData)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestGetERC721(t *testing.T) {
	var contracts []string = []string{
		"0x4cde5683f5f5616a8919a1d487552f2454c47a33",
		"0x828b5adbe8e0a06aaf5d5a5fd16e6b15a393d86e",
		"0xe70d5dd89ee10e55d88167b47169296143e3268b",
		"0x5e74094cd416f55179dbd0e45b1a8ed030e396a1",
		"0xea3dcb2d8c0fa2f50f6dfd1244749d0a6d0d9e13",
		"0x8a3acae3cd954fede1b46a6cc79f9189a6c79c56",
		"0x5bc94e9347f3b9be8415bdfd24af16666704e44f",
		"0x90e88d4c8e8f19af15dfeabd516d80666f06a2f5",
		"0xbe7095dbbe04e8374ea5f9f5b3f30a48d57cb004",
		"0xc014b45d680b5a4bf51ccda778a68d5251c14b5e",
		"0xd8d638be21b4101e1858de84d4540aed2d02674d",
		"0x443b862d3815b1898e85085cafca57fc4335a1be",
		"0x2aeaffc99cef9f6fc0869c1f16f890abdfcc222b",
		"0x4b8a1b8cf55d0af67c308ce22b4e4c11b04faeb0",
		"0x2a187453064356c898cae034eaed119e1663acb8",
		"0x27b4bc90fbe56f02ef50f2e2f79d7813aa8941a7",
		"0x57f0b53926dd62f2e26bc40b30140abea474da94",
		"0x57f1887a8bf19b14fc0df6fd9b2acc9af147ea85",
		"0xf85c874ea05e2225982b48c93a7c7f701065d91e",
		"0x6cc462bc49cecfe943bc4f477b23b92906e6074f",
		"0xc1caf0c19a8ac28c41fe59ba6c754e4b9bd54de9",
		"0x452b2bc7c94515720b36d304ce33909a8323f3e3",
		"0x3910d4afdf276a0dc8af632ccfceccf5ba04a3b7",
		"0x22c1f6050e56d2876009903609a2cc3fef83b415",
		"0xdc76a2de1861ea49e8b41a1de1e461085e8f369f",
	}

	for _, contract := range contracts {
		tokens, err := GetERC721(map[string]interface{}{"contract_addr": contract})
		if err != nil {
			t.Error(err)
			return
		}

		for i := 0; i < len(tokens); i++ {
			for j := 0; j < len(tokens); j++ {
				if i == j {
					continue
				}
				if tokens[i].TokenID == tokens[j].TokenID {
					fmt.Printf("重复的ID: %d 合约地址: %s \n", tokens[i].TokenID,contract)
				}
			}
		}
	}

}
