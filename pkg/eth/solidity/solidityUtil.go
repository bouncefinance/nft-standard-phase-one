package solidity

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
)

const (
	STRING  = "string"
	BOOL    = "bool"
	ADDRESS = "common.Address"
	UINT8   = "uint8"
	UINT16  = "uint16"
	UINT32  = "uint32"
	UINT64  = "uint64"
	UINT256 = "uint"
	INT8    = "int8"
	INT16   = "int16"
	INT32   = "int32"
	INT64   = "int64"
	INT256  = "int"
	BIGINT  = "*big.Int"
	BYTES   = "[]uint8" //	[]uint8 == []byte
)

var GoSolidityType map[string]string

func init() {
	GoSolidityType = make(map[string]string, 0)
	GoSolidityType[STRING] = STRING
	GoSolidityType[BOOL] = BOOL
	GoSolidityType[ADDRESS] = "address"
	GoSolidityType[UINT8] = UINT8
	GoSolidityType[UINT16] = UINT16
	GoSolidityType[UINT32] = UINT32
	GoSolidityType[UINT64] = UINT64
	GoSolidityType[UINT256] = "uint256"
	GoSolidityType[INT8] = INT8
	GoSolidityType[INT16] = INT16
	GoSolidityType[INT32] = INT32
	GoSolidityType[INT64] = INT64
	GoSolidityType[INT256] = "int256"
	GoSolidityType[BIGINT] = "uint256"
	GoSolidityType[BYTES] = "bytes"
}

func ABIEncode(values []interface{}) ([]byte, error) { //types []string,
	types := make([]abi.Type, len(values))
	for key, value := range values {
		typeStr := reflect.TypeOf(value).String()
		typeSoli := GoSolidityType[typeStr]
		ty, err := abi.NewType(typeSoli, "", nil)
		if err != nil {
			return nil, err
		}
		types[key] = ty

		switch typeStr {
		case UINT8:
			values[key] = big.NewInt(int64(value.(uint8)))
		case UINT16:
			values[key] = big.NewInt(int64(value.(uint16)))
		case UINT32:
			values[key] = big.NewInt(int64(value.(uint32)))
		case UINT64:
			values[key] = big.NewInt(int64(value.(uint64)))
		case UINT256:
			values[key] = big.NewInt(int64(value.(uint)))
		case INT8:
			values[key] = big.NewInt(int64(value.(int8)))
		case INT16:
			values[key] = big.NewInt(int64(value.(int16)))
		case INT32:
			values[key] = big.NewInt(int64(value.(int32)))
		case INT64:
			values[key] = big.NewInt(value.(int64))
		case INT256:
			values[key] = big.NewInt(int64(value.(int)))
		}
	}

	arg := abi.Arguments{}
	for _, typ := range types {
		arg = append(arg, abi.Argument{
			Type: typ,
		})
	}

	return arg.Pack(values...)
}


func ABIEncodeWithSignature(method string, values []interface{}) ([]byte, error) {
	methodByte := crypto.Keccak256([]byte(method))
	methodSign := methodByte[:4]

	argsData, err := ABIEncode(values)
	if err != nil {
		return nil, err
	}
	return append(methodSign, argsData...), nil
}

func ABIDecodeString(data []byte)(string,error) {
	stringTy, err := abi.NewType("string", "", nil)
	if err != nil {
		return "", errors.Wrap(err,"abi.NewType error")
	}
	args := abi.Arguments{
		{
			Type: stringTy,
		},
	}

	object := &struct {
		Str string
	}{}
	err = args.Unpack(object, data)
	if err != nil {
		return "", errors.Wrap(err,"args.Unpack error")
	}
	return object.Str,nil
}
