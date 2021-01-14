package solidity

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestABIEncode(t *testing.T) {
	i := big.NewInt(299999)
	//address := common.HexToAddress("0x7F34103Fb28086E4BE5b0b07CaEe5F7c4E58Fe9F")
	//ok := true

	//hexByte, err := hex.DecodeString("ffffff")
	//if err != nil {
	//	t.Error(err)
	//}
	bytes, err := ABIEncode([]interface{}{i})
	if err != nil {
		t.Error(err)
	}
	s := hex.EncodeToString(bytes)
	fmt.Println(s)
}

func TestABIDecodeString(t *testing.T) {
	str := ""
	bytes, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003668747470733a2f2f7368616c6f6d68752e6769746875622e696f2f626f756e63652f4e46542f746f6b656e55524c5f30332e6a736f6e00000000000000000000")
	if err != nil {
		t.Error(err)
		return
	}
	ABIDecodeString(bytes)
}
