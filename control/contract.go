package control

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
	"github.com/wupeaking/vechaintool/models"
	"github.com/wupeaking/vechaintool/vechainclient"
	"github.com/wupeaking/vechaintool/view"
	"math/big"
	"strconv"
	"time"
)

// ViewFuncCall 调用视图函数
func ViewFuncCall(method string, args []string, log *ui.MultilineEntry) {
	if method == "" {
		log.Append("method cannot be empty")
		return
	}

	abiObj := models.Setting.ABIObj()
	if abiObj == nil {
		log.Append("ABI file not set")
		return
	}

	if len(abiObj.Methods[method].Inputs) != len(args) {
		log.Append("wrong number of parameters")
		return
	}
	argsI := make([]interface{}, 0)

	for i, arg := range abiObj.Methods[method].Inputs {
		switch arg.Type.String() {
		case "uint8":
			fallthrough
		case "uint16":
			fallthrough
		case "uint32":
			fallthrough
		case "uint64":
			value, err := strconv.ParseUint(args[i], 10, 0)
			if err != nil {
				log.Append("string conversion to number failed: " + err.Error())
				return
			} else {
				argsI = append(argsI, value)
			}

		case "uint128":
			fallthrough
		case "uint256":
			fallthrough
		case "int128":
			fallthrough
		case "int256":
			value, ok := new(big.Int).SetString(args[i], 10)
			if !ok {
				log.Append("string conversion to number failed")
				return
			} else {
				argsI = append(argsI, value)
			}
		case "int8":
			fallthrough
		case "int16":
			fallthrough
		case "int32":
			fallthrough
		case "int64":
			value, err := strconv.ParseInt(args[i], 10, 0)
			if err != nil {
				log.Append("string conversion to number failed: " + err.Error())
				return
			} else {
				argsI = append(argsI, value)
			}

		case "bool":
			if args[i] == "true" {
				argsI = append(argsI, true)
			} else if args[i] == "false" {
				argsI = append(argsI, false)
			} else {
				log.Append(fmt.Sprintf("the args[%d] must be true/false\n", i))
			}

		case "address":
			if len(args[i]) != 42 {
				log.Append(fmt.Sprintf("the args[%d] must be ethereum address\n", i))
				return
			}

			b, err := hex.DecodeString(args[i][2:])
			if err != nil {
				log.Append(fmt.Sprintf("hexadecimal decoding failed %s\n", err.Error()))
				return
			}
			var addr [20]byte
			copy(addr[:], b)
			argsI = append(argsI, addr)

		default:
			argsI = append(argsI, args[i])
		}
	}

	// 进行编码
	inputByte, err := abiObj.Pack(method, argsI...)
	if err != nil {
		log.Append(fmt.Sprintf("ABI encoding failed: %s\n", err.Error()))
		return
	}

	//发起调用
	veCli, _ := vechain.NewVeChainClient(models.Setting.RPC, "", "", 10)

	ret, err := veCli.SimulateContract(context.Background(), fmt.Sprintf("0x%x", inputByte),
		"0", models.Setting.Contract)
	if err != nil {
		log.Append(fmt.Sprintf("call contract failed %s\n", err.Error()))
		return
	}

	if ret.Reverted {
		log.Append(fmt.Sprintf("contract execute reverted  %s\n", ret.VMErr))
		return
	}

	// 对返回的结果进行解包
	// 去除0x
	retData, err := hex.DecodeString(ret.Data[2:])
	if err != nil {
		log.Append(fmt.Sprintf("Failed to convert hex conversion to result %s\n", err.Error()))
		return
	}

	retI := make([]interface{}, 0)
	for _, arg := range abiObj.Methods[method].Outputs {
		switch arg.Type.String() {
		case "uint8":
			var tmp uint8
			retI = append(retI, &tmp)
		case "uint16":
			var tmp uint16
			retI = append(retI, &tmp)
		case "uint32":
			var tmp uint32
			retI = append(retI, &tmp)
		case "uint64":
			var tmp uint64
			retI = append(retI, &tmp)

		case "uint128":
			fallthrough
		case "uint256":
			fallthrough
		case "int128":
			fallthrough
		case "int256":
			tmp := big.NewInt(0)
			retI = append(retI, &tmp)

		case "int8":
			var tmp int8
			retI = append(retI, &tmp)
		case "int16":
			var tmp int16
			retI = append(retI, &tmp)
		case "int32":
			var tmp int32
			retI = append(retI, &tmp)
		case "int64":
			var tmp int64
			retI = append(retI, &tmp)

		case "bool":
			var tmp bool
			retI = append(retI, &tmp)

		case "address":
			var addr common.Address
			retI = append(retI, &addr)

		default:
			var tmp string
			retI = append(retI, &tmp)
		}
	}

	if len(retI) == 1 {
		ret := retI[0]
		err = abiObj.Unpack(ret, method, retData)
		if err != nil {
			log.Append(fmt.Sprintf("Decoding the call result failed %s\n", err.Error()))
			return
		}
		arg := abiObj.Methods[method].Outputs[0]
		data, _ := json.Marshal(ret)
		log.Append(fmt.Sprintf("args[%d](%s): %v \n", 0, arg.Name, string(data)))
		return
	}

	err = abiObj.Unpack(&retI, method, retData)
	if err != nil {
		log.Append(fmt.Sprintf("Decoding the call result failed %s\n", err.Error()))
		return
	}

	for i, arg := range abiObj.Methods[method].Outputs {
		d, _ := json.Marshal(retI[i])
		log.Append(fmt.Sprintf("args[%d](%s): %v \n", i, arg.Name, string(d)))
	}
	return
}

func CallContract(method string, args []string, log *ui.MultilineEntry) {
	if method == "" {
		log.Append("method cannot be empty")
		return
	}

	abiObj := models.Setting.ABIObj()
	if abiObj == nil {
		log.Append("ABI file not set\n")
		return
	}
	if len(abiObj.Methods[method].Inputs) != len(args) {
		log.Append("wrong number of parameters\n")
		return
	}
	if models.Setting.PrivateKey() == nil {
		log.Append("private key not set\n")
		return
	}
	pk := models.Setting.PrivateKey().D.Bytes()

	argsI := make([]interface{}, 0)

	for i, arg := range abiObj.Methods[method].Inputs {
		switch arg.Type.String() {
		case "uint8":
			fallthrough
		case "uint16":
			fallthrough
		case "uint32":
			fallthrough
		case "uint64":
			value, err := strconv.ParseUint(args[i], 10, 0)
			if err != nil {
				log.Append("string conversion to number failed\n")
				return
			} else {
				argsI = append(argsI, value)
			}

		case "uint128":
			fallthrough
		case "uint256":
			fallthrough
		case "int128":
			fallthrough
		case "int256":
			value, ok := new(big.Int).SetString(args[i], 10)
			if !ok {
				log.Append(fmt.Sprintf("string conversion to number failed\n"))
				return
			} else {
				argsI = append(argsI, value)
			}
		case "int8":
			fallthrough
		case "int16":
			fallthrough
		case "int32":
			fallthrough
		case "int64":
			value, err := strconv.ParseInt(args[i], 10, 0)
			if err != nil {
				log.Append(fmt.Sprintf("string conversion to number failed\n"))
				return
			} else {
				argsI = append(argsI, value)
			}

		case "bool":
			if args[i] == "true" {
				argsI = append(argsI, true)
			} else if args[i] == "false" {
				argsI = append(argsI, false)
			} else {
				log.Append(fmt.Sprintf("args[%d] must be true/false\n", i))
				return
			}

		case "address":
			if len(args[i]) != 42 {
				log.Append(fmt.Sprintf("args[%d] must be ethereum address \n", i))
				return
			}

			b, err := hex.DecodeString(args[i][2:])
			if err != nil {
				log.Append(fmt.Sprintf("Failed to convert hex conversion to result: %s\n", err.Error()))
				return
			}
			var addr [20]byte
			copy(addr[:], b)
			argsI = append(argsI, addr)

		default:
			argsI = append(argsI, args[i])
		}
	}

	// 进行编码
	inputByte, err := abiObj.Pack(method, argsI...)
	if err != nil {
		log.Append(fmt.Sprintf("Input encoding failed: %s\n", err.Error()))
		return
	}

	//发起调用
	veCli, _ := vechain.NewVeChainClient(models.Setting.RPC, "", "", 10)
	// 获取最新区块
	blk, err := veCli.BlockInfo(context.Background(), 0)
	if err != nil {
		log.Append(fmt.Sprintf("query blockchain failed: %s\n", err.Error()))
		return
	}
	// 开始构造交易
	// todo:: gas 暂时设置5000000
	var chain byte
	if false {
		chain = 0x4a
	} else {
		chain = 0x27
	}
	tid, raw, err := constructRawTx(models.Setting.Contract, pk, uint32(blk.BlockNum),
		new(big.Int).SetUint64(0), inputByte, 5000000, chain)

	if err != nil {
		log.Append(fmt.Sprintf("Build original transaction failed: %s\n", err.Error()))
		return
	}

	msg := "确认发起此交易?"
	view.ConfirmDialog(msg, func() {
		// 开始广播交易
		_, err = veCli.PushTx(context.Background(), raw)
		if err != nil {
			log.Append(fmt.Sprintf("Broadcast transaction failed: %s\n", err.Error()))
			return
		}
		log.Append(fmt.Sprintf("Broadcast transaction success, tx_id: %s \n", tid))
	}, func() {
		log.Append(fmt.Sprintf("cancel transaction\n"))
		return
	})
	return
}

func constructRawTx(to string, prvk []byte, blockNum uint32,
	amount *big.Int, input []byte, gas uint64, chain byte) (string, string, error) {
	toAddr, err := thor.ParseAddress(to)
	if err != nil {
		return "", "", err
	}
	//   chaintag  创世区块ID 最后一个字节 测试链为0x27 生产链为0x4a
	trx := new(tx.Builder).ChainTag(chain).
		BlockRef(tx.NewBlockRef(blockNum)).
		Expiration(720).
		Clause(tx.NewClause(&toAddr).WithValue(amount).WithData(input)).
		GasPriceCoef(0).
		Gas(gas).
		DependsOn(nil).
		Nonce(uint64(time.Now().UnixNano())).Build()

	priv, err := crypto.ToECDSA(prvk) //  HexToECDSA(b.Text(16))
	if err != nil {
		return "", "", err
	}
	sig, err := crypto.Sign(trx.SigningHash().Bytes(), priv)
	if err != nil {
		return "", "", err
	}

	trx = trx.WithSignature(sig)
	d, err := rlp.EncodeToBytes(trx)
	if err != nil {
		return "", "", err
	}
	return trx.ID().String(), hex.EncodeToString(d), nil
}
