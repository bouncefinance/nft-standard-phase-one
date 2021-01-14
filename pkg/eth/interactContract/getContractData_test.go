package interactContract

import (
	"fmt"
	"testing"
)

func TestGetAllTXFromContract(t *testing.T) {
	address := "0x5bc94e9347f3b9be8415bdfd24af16666704e44f"
	chainID := 56

	txes, err := GetAllNormalTX(address, chainID)
	if err!=nil{
		t.Error(err)
		return
	}

	fmt.Println(len(txes))
}

func TestGetAll721TX(t *testing.T) {
	address := "0xdf7952b35f24acf7fc0487d01c8d5690a60dba07"
	chainID := 56
	GetAll721TX(address,chainID)
}

func TestGetNormalTXWithBlockNum(t *testing.T) {
	address := "0x7f15017506978517Db9eb0Abd39d12D86B2Af395"
	chainID := 4
	from:="7437324"
	to:="7653821"
	txes, err := GetNormalTXWithBlockNum(address, chainID, from, to)
	if err!=nil{
		t.Error(err)
		return
	}
	fmt.Println(txes)
}