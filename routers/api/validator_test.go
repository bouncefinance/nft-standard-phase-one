package api

import (
	"Ankr-gin-ERC721/pkg/eth/interactContract"
	"testing"
)

func TestIsSupportInterface(t *testing.T) {
	contractAddr:="0x7f15017506978517Db9eb0Abd39d12D86B2Af395"
	chainID:=4

	IsSupportInterface(contractAddr,interactContract.ERC1155_INTERFACE_ID,chainID)
}
