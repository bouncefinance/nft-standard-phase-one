package interactContract

import (
	"Ankr-gin-ERC721/conf"
	"Ankr-gin-ERC721/pkg/eth/solidity"
	"Ankr-gin-ERC721/pkg/setting"
	"Ankr-gin-ERC721/pkg/util"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

const (
	OWNEROF_SIGN           = "0x6352211e"
	BALANCEOF_SIGN         = "0x70a08231"
	SUPPORT_INTERFACE_SIGN = "0x01ffc9a7"
	GET_CONFIG_SIGN        = "0x6dd5b69d"
	GET_CONFIG_SIGN01      = "0x8ec872e3"
	GET_CONFIG_ADDR_SIGN   = "0x52665f47"
	GET_TOKEN_URI          = "0xc87b56dd" //	keccak256(abi.encodePacked("tokenURI(uint256)"));
	GET_1155_BASE_URI      = "0x0e89341c" //	keccak256(abi.encodePacked("uri(uint256)"));

	JSONRPC = "2.0"
	METHOD  = "eth_call"
	BLOCK   = "latest"
)
const (
	//	ERC165
	ERC721_INTERFACE_ID  = "80ac58cd"
	ERC1155_INTERFACE_ID = "d9b67a26"
	ACTIVELY_RESPONSE = 1
)

var (
	requestID = 1
)

type RPCData struct {
	JsonRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      int    `json:"id"`
	Params  []interface{}
}
type Param struct {
	//From     string `json:"from"`
	To string `json:"to"`
	//Gas      uint   `json:"gas"`
	//GasPrice uint   `json:"gasPrice"`
	//Value    uint   `json:"value"`
	Data string `json:"data"`
}
type ReturnData struct {
	JsonRPC string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}
type EtherscanData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

//	https://eth.wiki/json-rpc/API eth_call接口

func SupportInterface(contractAddr string, interSign string, chainID int) ([]byte, error) {
	bytes4Ty, err := abi.NewType("bytes4", "", nil)
	if err != nil {
		return nil, errors.Wrap(err,"abi.NewType error")
	}

	arg := abi.Arguments{
		{
			Type: bytes4Ty,
		},
	}

	var interSignBytes4 [4]byte
	interSignBytes, _ := hex.DecodeString(interSign)
	for i := 0; i < len(interSignBytes4); i++ {
		interSignBytes4[i] = interSignBytes[i]
	}

	bytes, err := arg.Pack(interSignBytes4)
	if err != nil {
		return nil, errors.Wrap(err,"arg.Pack error")
	}
	hex_ := hex.EncodeToString(bytes)
	param := Param{
		To:   contractAddr,
		Data: SUPPORT_INTERFACE_SIGN + string(hex_),
	}

	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  METHOD,
		Id:      requestID,
		Params:  make([]interface{}, 0),
	}
	rpcData.Params = append(rpcData.Params, param)
	rpcData.Params = append(rpcData.Params, BLOCK)

	url:=""
	if chainID == conf.BSC_CHAINID || chainID == conf.BSC_TEST_CHAINID{
		url = conf.ETHHttpURLs[chainID]
	}else {
		url = conf.ETHHttpURLs[chainID] + conf.ConfigMsg.ProjectID
	}
	data, err := util.PostUrl(url, nil, rpcData, map[string]string{"Content-type":"application/json"})
	if err != nil {
		return data, errors.Wrap(err,"util.PostUrl error")
	}
	requestID++

	return data, nil
}

func BalanceOf(contractAddr string, userAddr string, chainID int) ([]byte, error) {
	bytes, err := solidity.ABIEncode([]interface{}{userAddr})
	if err != nil {
		return nil, err
	}
	hex_ := hex.EncodeToString(bytes)

	param := Param{
		To:   contractAddr,
		Data: BALANCEOF_SIGN + string(hex_),
	}
	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  METHOD,
		Id:      requestID,
		Params:  make([]interface{}, 0),
	}
	rpcData.Params = append(rpcData.Params, param)
	rpcData.Params = append(rpcData.Params, BLOCK)

	url:=""
	if chainID == conf.BSC_CHAINID || chainID == conf.BSC_TEST_CHAINID{
		url = conf.ETHHttpURLs[chainID]
	}else {
		url = conf.ETHHttpURLs[chainID] + conf.ConfigMsg.ProjectID
	}
	data, err := util.PostUrl(url, nil, rpcData, nil)
	if err != nil {
		return data, err
	}
	requestID++

	return data, nil
}

func OwnerOf(tokenID int, contractAddr string, chainID int) ([]byte, error) {
	bytes, err := solidity.ABIEncode([]interface{}{tokenID})
	if err != nil {
		return nil, err
	}
	hex := hex.EncodeToString(bytes)

	param := Param{
		To:   contractAddr,
		Data: OWNEROF_SIGN + hex,
	}
	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  METHOD,
		Id:      requestID,
		Params:  make([]interface{}, 0),
	}
	rpcData.Params = append(rpcData.Params, param)
	rpcData.Params = append(rpcData.Params, BLOCK)

	url:=""
	if chainID == conf.BSC_CHAINID || chainID == conf.BSC_TEST_CHAINID{
		url = conf.ETHHttpURLs[chainID]
	}else {
		url = conf.ETHHttpURLs[chainID] + conf.ConfigMsg.ProjectID
	}

	data, err := util.PostUrl(url, nil, rpcData, nil)
	if err != nil {
		return data, err
	}
	requestID++

	return data, nil
}

func GetConfig(contractAddr string, chainID int, key [32]byte) ([]byte, error) {
	bytes32Ty, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments{
		{
			Type: bytes32Ty,
		},
	}
	bytes, err := args.Pack(key)
	if err != nil {
		return nil, err
	}

	hexS := hex.EncodeToString(bytes)

	param := Param{
		To:   contractAddr,
		Data: GET_CONFIG_SIGN + hexS,
	}
	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  METHOD,
		Id:      requestID,
		Params:  make([]interface{}, 0),
	}
	rpcData.Params = append(rpcData.Params, param)
	rpcData.Params = append(rpcData.Params, BLOCK)

	url:=""
	if chainID == conf.BSC_CHAINID || chainID == conf.BSC_TEST_CHAINID{
		url = conf.ETHHttpURLs[chainID]
	}else {
		url = conf.ETHHttpURLs[chainID] + conf.ConfigMsg.ProjectID
	}

	data, err := util.PostUrl(url, nil, rpcData, nil)
	if err != nil {
		return data, err
	}
	requestID++

	return data, nil
}

func GetConfig01(contractAddr string, chainID int, pre [32]byte, index uint) ([]byte, error) {
	bytes32Ty, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, err
	}
	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments{
		{
			Type: bytes32Ty,
		},
		{
			Type: uint256Ty,
		},
	}
	bytes, err := args.Pack(pre, big.NewInt(int64(index)))
	if err != nil {
		return nil, err
	}

	hex := hex.EncodeToString(bytes)

	param := Param{
		To:   contractAddr,
		Data: GET_CONFIG_SIGN01 + hex,
	}
	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  METHOD,
		Id:      requestID,
		Params:  make([]interface{}, 0),
	}
	rpcData.Params = append(rpcData.Params, param)
	rpcData.Params = append(rpcData.Params, BLOCK)

	url:=""
	if chainID == conf.BSC_CHAINID || chainID == conf.BSC_TEST_CHAINID{
		url = conf.ETHHttpURLs[chainID]
	}else {
		url = conf.ETHHttpURLs[chainID] + conf.ConfigMsg.ProjectID
	}

	data, err := util.PostUrl(url, nil, rpcData, nil)
	if err != nil {
		return data, err
	}
	requestID++

	return data, nil
}

func GetConfigWithAddr(contractAddr string, chainID int, pre [32]byte, addr common.Address) ([]byte, error) {
	bytes32Ty, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, err
	}
	addressTy, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments{
		{
			Type: bytes32Ty,
		},
		{
			Type: addressTy,
		},
	}
	bytes, err := args.Pack(pre, addr)
	if err != nil {
		return nil, err
	}

	hex := hex.EncodeToString(bytes)

	param := Param{
		To:   contractAddr,
		Data: GET_CONFIG_ADDR_SIGN + hex,
	}
	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  METHOD,
		Id:      requestID,
		Params:  make([]interface{}, 0),
	}
	rpcData.Params = append(rpcData.Params, param)
	rpcData.Params = append(rpcData.Params, BLOCK)

	url:=""
	if chainID == conf.BSC_CHAINID || chainID == conf.BSC_TEST_CHAINID{
		url = conf.ETHHttpURLs[chainID]
	}else {
		url = conf.ETHHttpURLs[chainID] + conf.ConfigMsg.ProjectID
	}

	data, err := util.PostUrl(url, nil, rpcData, nil)
	if err != nil {
		return data, err
	}
	requestID++

	return data, nil
}

func GetContractABI(contractAddr string) string {
	//https://api.etherscan.io/api?module=contract&action=getabi&address=0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413&apikey=YourApiKeyToken
	//url := "https://api-rinkeby.etherscan.io/api"
	params := make(map[string]string)
	params["module"] = "contract"
	params["action"] = "getabi"
	params["address"] = contractAddr
	params["apikey"] = conf.ConfigMsg.ApiKeyToken

	data, err := util.PostUrl(setting.RINKEBY_ETHERSCAN_URL, params, nil, nil)
	if err != nil {
		fmt.Println("GetContractABI error: ", err)
		return ""
	}
	etherscanData := EtherscanData{}
	err = json.Unmarshal(data, &etherscanData)
	if err != nil {
		fmt.Println("GetContractABI json.Unmarshal error: ", err)
		return ""
	}

	return etherscanData.Result
}

func GetTokenURI(contractAddr string, chainID int, tokenID int) ([]byte, error) {
	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, errors.Wrap(err,"abi.NewType error")
	}
	args := abi.Arguments{
		{
			Type: uint256Ty,
		},
	}

	bytes, err := args.Pack(big.NewInt(int64(tokenID)))
	if err != nil {
		return nil, errors.Wrap(err,"args.Pack error")
	}
	hex := hex.EncodeToString(bytes)

	param := Param{
		To:   contractAddr,
		Data: GET_TOKEN_URI + hex,
	}

	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  METHOD,
		Id:      requestID,
		Params:  make([]interface{}, 0),
	}
	rpcData.Params = append(rpcData.Params, param)
	rpcData.Params = append(rpcData.Params, BLOCK)

	url:=""
	if chainID == conf.BSC_CHAINID || chainID == conf.BSC_TEST_CHAINID{
		url = conf.ETHHttpURLs[chainID]
	}else {
		url = conf.ETHHttpURLs[chainID] + conf.ConfigMsg.ProjectID
	}

	time:=0
	data, err := util.PostUrl(url, nil, rpcData, map[string]string{"Content-type":"application/json"})
	for err != nil && time < conf.RETRY_TIME {
		data, err = util.PostUrl(url, nil, rpcData, map[string]string{"Content-type":"application/json"})
		time++
	}
	if err != nil {
		return data, errors.Wrap(err,"util.PostUrl error")
	}
	requestID++

	return data, nil
}

func GetBlockNum()([]byte, error) {
	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  "eth_blockNumber",
		Id:      requestID,
		Params:  nil,
	}
	url := conf.ETHHttpURLs[conf.MAIN_CHAINID]

	data, err := util.PostUrl(url, nil, rpcData, map[string]string{"Content-type":"application/json"})
	if err!=nil{
		return nil,errors.Wrap(err,"PostURL error")
	}
	fmt.Println(string(data))

	return data,nil
}

func GetTokenURIStr(contractAddr string, chainID int, tokenID int) (string, error) {
	uriData, err := GetTokenURI(contractAddr, chainID, tokenID)
	if err != nil {
		return "", errors.Wrap(err,"GetTokenURI error")
	}

	result := ReturnData{}
	err = json.Unmarshal(uriData, &result)
	if err != nil {
		return "", errors.Wrap(err,"json.Unmarshal error")
	}
	if len(result.Result) == 0{
		return "", nil
	}
	bytes, err := hex.DecodeString(result.Result[2:])
	if err != nil {
		return "", errors.Wrap(err,"hex.DecodeString error")
	}

	decodeString, err := solidity.ABIDecodeString(bytes)
	if err != nil {
		return "", errors.Wrap(err,"solidity.ABIDecodeString error")
	}
	return decodeString, nil
}

func Get1155BaseURI(contractAddr string, chainID int) (string, error) {
	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return "", err
	}
	args := abi.Arguments{
		{
			Type: uint256Ty,
		},
	}

	bytes, err := args.Pack(big.NewInt(int64(1)))
	if err != nil {
		return "", err
	}
	hexStr := hex.EncodeToString(bytes)

	param := Param{
		To:   contractAddr,
		Data: GET_1155_BASE_URI + hexStr,
	}

	rpcData := RPCData{
		JsonRPC: JSONRPC,
		Method:  METHOD,
		Id:      requestID,
		Params:  make([]interface{}, 0),
	}
	rpcData.Params = append(rpcData.Params, param)
	rpcData.Params = append(rpcData.Params, BLOCK)

	url:=""
	if chainID == conf.BSC_CHAINID || chainID == conf.BSC_TEST_CHAINID{
		url = conf.ETHHttpURLs[chainID]
	}else {
		url = conf.ETHHttpURLs[chainID] + conf.ConfigMsg.ProjectID
	}

	data, err := util.PostUrl(url, nil, rpcData, map[string]string{"Content-type":"application/json"})
	if err != nil {
		return "", err
	}
	requestID++

	result := ReturnData{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}
	if len(result.Result) == 0{
		return "", nil
	}
	bytes, err = hex.DecodeString(result.Result[2:])
	if err != nil {
		return "", err
	}

	decodeString, err := solidity.ABIDecodeString(bytes)
	if err != nil {
		return "", err
	}
	return decodeString, nil
}
