package interactContract

import (
	"Ankr-gin-ERC721/pkg/eth/solidity"
	"Ankr-gin-ERC721/pkg/util"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func TestBalanceOf(t *testing.T) {
	data, e := BalanceOf("0x93e508f373690cC4307a7A2363e573E63dAEF54E", "0xBCcC2073ADfC46421308f62cfD9868dF00D339a8", 4)
	if e != nil {
		t.Error("BalanceOf error: ", e)
	}
	fmt.Println("return data: ", string(data))
	i := util.BytesToInt(data)
	fmt.Println("--->", i)
}

func TestOwnerOf(t *testing.T) {
	data, e := OwnerOf(1, "0xf6C3Aa70f29B64BA74dd6Abe6728cf8e190011b5", 97)
	if e != nil {
		t.Error("OwnerOf error: ", e)
	}
	fmt.Println("return data: ", string(data))

	returnData := ReturnData{}
	err := json.Unmarshal(data, &returnData)
	if err != nil {
		t.Error("OwnerOf error: ", err)
	}

	splits := strings.Split(returnData.Result, "0x000000000000000000000000")
	fmt.Println(splits)
}

func TestSupportInterface(t *testing.T) {
	data, err := SupportInterface("0xDf7952B35f24aCF7fC0487D01c8d5690a60DBa07", ERC721_INTERFACE_ID, 56)
	if err != nil {
		t.Error("BalanceOf error: ", err)
	}

	fmt.Println(string(data))

	returnData := ReturnData{}
	err = json.Unmarshal(data, &returnData)
	if err != nil {
		t.Error("json.Unmarshal error: ", err)
	}

	i, err := strconv.ParseInt(returnData.Result, 0, 0)
	if err != nil {
		t.Error("strconv.ParseInt error: ", err)
	}

	fmt.Println("result ==>", i)
}

func TestGetConfig(t *testing.T) {
	/*pre := []byte("proposes")
	key := []byte{0x0}
	pre32 := common.LeftPadBytes(pre, 32)
	key32 := common.LeftPadBytes(key, 32)

	s := hex.EncodeToString(pre32)
	s1 := hex.EncodeToString(pre)
	fmt.Println(s,s1)

	var xor = [32]byte{}

	for i, v := range pre32 {
		xor[i] = v ^ key32[i]
	}

	xs := hex.EncodeToString(xor[:])
	fmt.Println(xs)*/

	//contractAddr := "0xa77A9FcbA2Ae5599e0054369d1655D186020ECE1"
	//chainID := 4
	contractAddr := "0x98945BC69A554F8b129b09aC8AfDc2cc2431c48E"
	chainID := 1

	bytes2Hex := common.Bytes2Hex([]byte(contractAddr))
	sss, _ := hex.DecodeString(bytes2Hex)
	fmt.Println("byte to hexString", sss)

	//pre :=  []byte("proposes")
	pre := crypto.Keccak256([]byte("proposes"))
	preH := crypto.Keccak256Hash([]byte("proposes"))
	fmt.Println("pre hash => ", preH.String())
	preByte := common.RightPadBytes(pre, 32)
	var pre32 [32]byte
	for i, v := range preByte {
		pre32[i] = v
	}

	ba, e := GetConfig01(contractAddr, chainID, pre32, 0)
	if e != nil {
		t.Error("GetConfig ERROR: ", e)
		return
	}
	fmt.Printf("result => %s \n\n", string(ba))

	result := ReturnData{}
	err := json.Unmarshal(ba, &result)
	if err != nil {
		t.Error(err)
	}
	bytes, err := hex.DecodeString(result.Result[2:])
	if err != nil {
		t.Error(err)
		return
	}

	i := 0
	xor := [32]byte{}
	for {
		lastPID, err := hex.DecodeString(result.Result[2:])
		if err != nil {
			t.Error(err)
		}
		for i, v := range pre32 {
			xor[i] = v ^ lastPID[i]
		}
		data, err := GetConfig(contractAddr, chainID, xor)
		if err != nil {
			t.Error("GetConfig ERROR: ", err)
		}
		fmt.Printf("result => %s \n\n", string(data))
		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Error(err)
		}
		i++
		if result.Result == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			if i == 3 {
				pre := []byte("proposeStatus")
				preByte := common.RightPadBytes(pre, 32)
				for i, v := range preByte {
					pre32[i] = v
				}

				continue
			}
			break
		}
		bytes, err := hex.DecodeString(result.Result[2:])
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println("查询的proposalID => ", string(bytes))
		temData := make([]byte, len(bytes))
		for i, byte := range bytes {
			temData[len(temData)-1-i] = byte
		}
		fmt.Println("查询的proposalID => OOOO ", string(bytes))
	}
	fmt.Println(i)
}

//	获得最新提案ID存储的数据
func TestGetConfig012(t *testing.T) {
	contractAddr := "0x98945BC69A554F8b129b09aC8AfDc2cc2431c48E"
	chainID := 1
	startProposalID := "bounce1"

	pre := []byte("proposer")
	//pre := []byte("proposeStatus")
	preByte := common.RightPadBytes(pre, 32)

	staProIDByte_ := []byte(startProposalID)
	staProIDByte := common.RightPadBytes(staProIDByte_, 32)

	var xor = [32]byte{}
	leftKey := staProIDByte
	for {
		fmt.Println("+++ >>>> ", hex.EncodeToString(leftKey))
		newInt := big.NewInt(0)
		//newByte := newInt.SetBytes(leftKey)
		//fmt.Println(newByte.String())

		newInt.SetString(hex.EncodeToString(leftKey), 16)
		fmt.Println("leftKey int => ", newInt.String())

		for i := 0; i < len(preByte); i++ {
			xor[i] = preByte[i] ^ leftKey[i]
		}
		data, err := GetConfig(contractAddr, chainID, xor)
		if err != nil {
			t.Error("GetConfig ERROR: ", err)
			return
		}
		result := ReturnData{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println("result => ", result)
		if result.Result == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			break
		}

		bytes, err := hex.DecodeString(result.Result[2:])
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println("查询的proposalID => ", string(bytes))

		leftKey = bytes
	}
}

func TestGetConfig01(t *testing.T) {
	pre := []byte("proposes")
	pre1 := []byte("shalom")
	bounce1 := []byte("bounce12")

	// value => 100
	pre32 := common.LeftPadBytes(pre, 32)
	pre13 := common.LeftPadBytes(pre1, 32)
	bounce1B := common.RightPadBytes(bounce1, 32)
	s := hex.EncodeToString(pre32)
	s13 := hex.EncodeToString(pre13)
	bounce1S := hex.EncodeToString(bounce1B)

	fmt.Println(s)
	fmt.Println(s13)
	fmt.Println(bounce1S)
}

func TestGetConfigWithAddr(t *testing.T) {
	contractAddr := "0x98945BC69A554F8b129b09aC8AfDc2cc2431c48E"
	chainID := 1
	pre := []byte("govRewardPerDay")
	preByte := common.RightPadBytes(pre, 32)
	var pre32 [32]byte
	for i, v := range preByte {
		pre32[i] = v
	}

	address := common.BigToAddress(big.NewInt(0))
	data, err := GetConfigWithAddr(contractAddr, chainID, pre32, address)
	if err != nil {
		fmt.Println("error: ", err)
		t.Error("error: ", err)
		return
	}
	fmt.Println("=== ", string(data))
}

func TestGovernanceID(t *testing.T) {
	//contractAddr := "0xa77A9FcbA2Ae5599e0054369d1655D186020ECE1"
	//chainID := 4

	contractAddr := "0x98945BC69A554F8b129b09aC8AfDc2cc2431c48E"
	chainID := 1
	proposerAddr := "0x0aff6665bb45bf349489b20e225a6c5d78e2280f"
	startProposalID := "0x06afb9c1f665a07f556d7807d7bb986bf3fb4d40ae2751e472a7f08e7c4cd43e"

	proposeStatus := crypto.Keccak256([]byte("proposeStatus"))
	proposer := crypto.Keccak256([]byte("proposer"))
	timePropose := crypto.Keccak256([]byte("timePropose"))
	proposeKey := crypto.Keccak256([]byte("proposeKey"))
	proposeValue := crypto.Keccak256([]byte("proposeValue"))
	proposes := crypto.Keccak256([]byte("proposes"))
	proposesVoting := crypto.Keccak256([]byte("proposesVoting"))
	votes := crypto.Keccak256([]byte("votes"))
	proposeContent := crypto.Keccak256([]byte("proposeContent"))
	voteYes := []byte("VOTE_YES")

	proposeStatusByte := common.RightPadBytes(proposeStatus, 32)
	proposerByte := common.RightPadBytes(proposer, 32)
	timeProposeByte := common.RightPadBytes(timePropose, 32)
	proposeKeyByte := common.RightPadBytes(proposeKey, 32)
	proposeValueByte := common.RightPadBytes(proposeValue, 32)
	proposesByte := common.RightPadBytes(proposes, 32)
	proposesVotingByte := common.RightPadBytes(proposesVoting, 32)
	votesByte := common.RightPadBytes(votes, 32)
	proposeContentByte := common.RightPadBytes(proposeContent, 32)
	voteYesByte := common.RightPadBytes(voteYes, 32)

	proIDData, err := hex.DecodeString(startProposalID[2:])
	if err != nil {
		t.Error("hex.DecodeString error: ", err)
		return
	}
	staProIDByte := common.RightPadBytes(proIDData, 32)

	var proposeStatusXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		proposeStatusXor[i] = proposeStatusByte[i] ^ staProIDByte[i]
	}
	var proposerXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		proposerXor[i] = proposerByte[i] ^ staProIDByte[i]
	}
	var timeProposeXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		timeProposeXor[i] = timeProposeByte[i] ^ staProIDByte[i]
	}
	var proposeKeyXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		proposeKeyXor[i] = proposeKeyByte[i] ^ staProIDByte[i]
	}
	var proposeValueXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		proposeValueXor[i] = proposeValueByte[i] ^ staProIDByte[i]
	}
	var proposesXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		proposesXor[i] = proposesByte[i] ^ staProIDByte[i]
	}
	var proposesVotingXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		proposesVotingXor[i] = proposesVotingByte[i] ^ staProIDByte[i]
	}
	proposerAddrByte := common.HexToAddress(proposerAddr).Bytes()
	proposerAddrByte = common.RightPadBytes(proposerAddrByte, 32)
	var votesYesSalt = [32]byte{}
	var votesXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		votesYesSalt[i] = staProIDByte[i] ^ voteYesByte[i]
	}
	for i := 0; i < len(staProIDByte); i++ {
		votesXor[i] = votesByte[i] ^ votesYesSalt[i]
	}
	var proposeContentXor = [32]byte{}
	for i := 0; i < len(staProIDByte); i++ {
		proposeContentXor[i] = proposeContentByte[i] ^ staProIDByte[i]
	}

	data, err := GetConfig(contractAddr, chainID, votesXor)
	result := ReturnData{}
	err = json.Unmarshal(data, &result)
	fmt.Println("votes => ", result)
	votesData, err := hex.DecodeString(result.Result[2:])
	newInt := big.NewInt(0)
	newInt = newInt.SetBytes(votesData)
	fmt.Println("votes => ", newInt.String())

	data, err = GetConfig(contractAddr, chainID, proposesVotingXor)
	err = json.Unmarshal(data, &result)
	fmt.Println("proposesVotingXor => ", result)

	data, err = GetConfig(contractAddr, chainID, proposesXor)
	err = json.Unmarshal(data, &result)
	fmt.Println("proposesXor => ", result)

	data, err = GetConfig(contractAddr, chainID, proposeValueXor)
	err = json.Unmarshal(data, &result)
	fmt.Println("proposeValueXor => ", result)

	data, err = GetConfig(contractAddr, chainID, proposeKeyXor)
	err = json.Unmarshal(data, &result)
	fmt.Println("proposeKeyXor => ", result)

	data, err = GetConfig(contractAddr, chainID, timeProposeXor)
	err = json.Unmarshal(data, &result)
	fmt.Println("timeProposeXor => ", result)

	data, err = GetConfig(contractAddr, chainID, proposerXor)
	err = json.Unmarshal(data, &result)
	fmt.Println("proposerXor => ", result)

	data, err = GetConfig(contractAddr, chainID, proposeContentXor)
	err = json.Unmarshal(data, &result)
	fmt.Println("proposeContent => ", result)

	data, err = GetConfig(contractAddr, chainID, proposeStatusXor)
	if err != nil {
		t.Error("GetConfig ERROR: ", err)
		return
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Error(err)
		return
	}
	//fmt.Println("proposeStatus => ", result)
	proposeStatusData, err := hex.DecodeString(result.Result[2:])
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("proposeStatus => ", string(proposeStatusData))

}

func TestGetTokenURI(t *testing.T) {
	contractAddr := "0xf6C3Aa70f29B64BA74dd6Abe6728cf8e190011b5"
	chainID := 1
	tokenID := 2
	data, err := GetTokenURI(contractAddr, chainID, tokenID)
	if err != nil {
		t.Error("GetConfig ERROR: ", err)
		return
	}
	fmt.Println(string(data))
	result := ReturnData{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(result)
	bytes, err := hex.DecodeString(result.Result[2:])
	if err != nil {
		t.Error(err)
		return
	}

	decodeString, err := solidity.ABIDecodeString(bytes)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("result => ", decodeString)
}

func TestGet1155BaseURI(t *testing.T) {
	contractAddr := "0x7f15017506978517Db9eb0Abd39d12D86B2Af395"
	chainID := 4
	uri, err := Get1155BaseURI(contractAddr, chainID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("uri => ", uri)
}

func TestGetTokenURIStr(t *testing.T) {
	contractAddr := "0x4cde5683f5f5616a8919a1d487552f2454c47a33"
	chainID := 56
	uri, err := GetTokenURIStr(contractAddr, chainID, 11411)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(uri)
}

func TestGetBlockNum(t *testing.T) {
	data, err := GetBlockNum()
	result := ReturnData{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("result => ", result)

	newInt := big.NewInt(0)
	newInt.SetString(result.Result, 0)
	fmt.Println("===>", newInt.Uint64())

	/*bytes, err := hex.DecodeString(result.Result[2:])
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("查询的proposalID => ", string(bytes))*/
}
