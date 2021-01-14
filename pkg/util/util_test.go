package util

import (
	"fmt"
	"testing"
)

func TestIntToBytes(t *testing.T) {
	intn, err := BytesToIntU([]byte("proposer"))
	if err != nil {
		t.Error("BytesToIntU ERROR:", err)
	}
	fmt.Println(intn)
}

func TestFileExist(t *testing.T) {
	exist := FileNotExist("/aaa")
	fmt.Printf("%t",exist)
}

func TestStrToLow(t *testing.T) {
	low := StrToLow("0xf942Bca10d7553867980d91A2c6428F0BD2b83B4")
	fmt.Println(low)
}