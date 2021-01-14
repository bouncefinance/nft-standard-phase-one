package runtime

import (
	"fmt"
	"testing"
)

func TestGetContractABI(t *testing.T) {
	abi := GetContractABI("0xf6C3Aa70f29B64BA74dd6Abe6728cf8e190011b5")
	fmt.Println("abi ==>",abi)
}

func TestParseLog(t *testing.T) {
	//ParseLog("0xf6C3Aa70f29B64BA74dd6Abe6728cf8e190011b5",)
}