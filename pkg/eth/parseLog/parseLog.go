package parseLog

import (
	"Ankr-gin-ERC721/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strconv"
	"strings"
)

const (
	MISSED = -1
	TRANSFERSINGLE = 0
	TRANSFERBATCH  = 1
	URI            = 2
)

var(
	SignSingle = crypto.Keccak256Hash([]byte("TransferSingle(address,address,address,uint256,uint256)"))
	SignBatch = crypto.Keccak256Hash([]byte("TransferBatch(address,address,address,uint256[],uint256[])"))
	SignURI = crypto.Keccak256Hash([]byte("URI(string,uint256)"))
	SignTransfer721 = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
)

type Batch struct {
	IDs    []*big.Int	`abi:"IDs"`
	Values []*big.Int	`abi:"Values"`
}


func ParseLog1155(log types.Log) (flag int, from string, to string,tokenIDs []uint,amount []*big.Int,uri string,tokenID int) {
	var (
		topics [4]string
	)
	for i := range log.Topics {
		topics[i] = log.Topics[i].Hex()
	}

	if topics[0] != SignSingle.Hex() && topics[0] != SignBatch.Hex() && topics[0] != SignURI.Hex() {
		flag = MISSED
		return
	}
	if topics[0] == SignSingle.Hex() || topics[0] == SignBatch.Hex(){
		from = "0x" + strings.Split(topics[2], "0x000000000000000000000000")[1]
		to = "0x" + strings.Split(topics[3], "0x000000000000000000000000")[1]
	}

	if topics[0] == SignSingle.Hex() {
		//event TransferSingle(address indexed _operator, address indexed _from, address indexed _to, uint256 _id, uint256 _value);
		uint256Ty01, _ := abi.NewType("uint256", "", nil)
		uint256Ty02, _ := abi.NewType("uint256", "", nil)

		arguments := abi.Arguments{
			{
				Type: uint256Ty01,
			},
			{
				Type: uint256Ty02,
			},
		}
		tem := make([]*big.Int, 2)
		err := arguments.Unpack(&tem, log.Data)
		if err != nil {
			flag = MISSED
			logger.Logger.Error().Str("contractAddress",log.Address.String()).Str("event name","TransferSingle").Str("from",from).Str("to",to).Msgf("ParseLog1155 arguments.Unpack error: %s",err)
			return
		}
		tokenIDs = append(tokenIDs, uint(tem[0].Uint64()))
		amount = append(amount, tem[1])
		flag = TRANSFERSINGLE
	} else if topics[0] == SignBatch.Hex() {
		//event TransferBatch(address indexed _operator, address indexed _from, address indexed _to, uint256[] _ids, uint256[] _values);
		uint256Ty01, _ := abi.NewType("uint256[]", "", nil)
		uint256Ty02, _ := abi.NewType("uint256[]", "", nil)

		arguments := abi.Arguments{
			{
				Type: uint256Ty01,
				Name: "IDs",
			},
			{
				Type: uint256Ty02,
				Name: "Values",
			},
		}
		tem := Batch{}
		err := arguments.Unpack(&tem, log.Data)
		if err != nil {
			flag = MISSED
			logger.Logger.Error().Str("contractAddress",log.Address.String()).Str("event name","TransferBatch").Str("from",from).Str("to",to).Msgf("ParseLog1155 arguments.Unpack error: %s",err)
			return
		}
		for _, idBig := range tem.IDs {
			tokenIDs = append(tokenIDs, uint(idBig.Uint64()))
		}
		//tokenIDs = tem.IDs
		amount = tem.Values
		flag = TRANSFERBATCH
	} else if topics[0] == SignURI.Hex() {
		//event URI(string _value, uint256 indexed _id);
		stringTy, _ := abi.NewType("string", "", nil)
		arguments := abi.Arguments{
			{
				Type: stringTy,
			},
		}
		tem := struct {
			URI string
		}{}
		err := arguments.Unpack(&tem, log.Data)
		if err != nil {
			flag = MISSED
			logger.Logger.Error().Str("contractAddress",log.Address.String()).Str("event name","URI").Msgf("ParseLog1155 arguments.Unpack error: %s",err)
			return
		}
		uri = tem.URI

		newInt := big.NewInt(0)
		tokenIDInt, ok := newInt.SetString(topics[1], 0)
		if !ok{
			flag = MISSED
			logger.Logger.Error().Str("contractAddress",log.Address.String()).Str("event name","URI").Str("tokenID topic",topics[1]).Msg("ParseLog1155 newInt.SetString error")
			return
		}
		tokenID = int(tokenIDInt.Int64())
		flag = URI
	}

	return
}

func ParseLog721(log types.Log) (bool, string, string, int) {
	var topics [4]string
	var (
		from    string
		to      string
		tokenID int
	)
	for i := range log.Topics {
		topics[i] = log.Topics[i].Hex()
	}
	if topics[0] != SignTransfer721.Hex() {
		return false, "", "", 0
	}

	from = "0x" + strings.Split(topics[1], "0x000000000000000000000000")[1]
	to = "0x" + strings.Split(topics[2], "0x000000000000000000000000")[1]
	tem, _ := strconv.ParseInt(topics[3], 0, 0)
	tokenID = int(tem)

	return true, from, to, tokenID
}
