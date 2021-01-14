package models

import (
	"testing"
)

func TestGetERC1155(t *testing.T) {
	erc1155 := ERC1155{
		ContractAddr: "0x7f15017506978517Db9eb0Abd39d12D86B2Af395",
		TokenID:      0,
		OwnerAddr:    "0xBCcC2073ADfC46421308f62cfD9868dF00D339a8",
		Balance:      "",
	}

	err := AddERC1155([]ERC1155{erc1155})
	if err!=nil{
		t.Error("AddERC1155 err: ",err)
	}
}
