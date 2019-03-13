package control

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/ethereum/go-ethereum/common"
	"github.com/wupeaking/vechaintool/models"
	"github.com/wupeaking/vechaintool/vechainclient"
	"math/big"
	"strconv"
)

// ViewFuncCall 调用视图函数
func ViewFuncCall(method string, args []string, log *ui.MultilineEntry  ) {
	if method == "" {
		log.Append("方法名称不能为空")
		return
	}

	abiObj := models.Setting.ABIObj()
	if abiObj == nil {
		log.Append("未设置ABI文件")
		return
	}


	if len(abiObj.Methods[method].Inputs) != len(args) {
		log.Append("参数个数错误")
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
				log.Append("字符串转换成数字失败: "+err.Error())
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
				 log.Append("字符串转换成数字失败")
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
				log.Append("字符串转换成数字失败: "+err.Error())
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
				log.Append(fmt.Sprintf("第%d个参数必须是true/false\n", i))
			}

		case "address":
			if len(args[i]) != 42 {
				log.Append(fmt.Sprintf("第%d个参数必须是以太坊地址\n", i))
				return
			}

			b, err := hex.DecodeString(args[i][2:])
			if err != nil {
				log.Append(fmt.Sprintf("十六进制解码失败 %s\n", err.Error()))
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
		log.Append(fmt.Sprintf("ABI编码失败 %s\n", err.Error()))
		return
	}

	//发起调用
	veCli, _ := vechain.NewVeChainClient(models.Setting.RPC, "", "", 10)

	ret, err := veCli.SimulateContract(context.Background(), fmt.Sprintf("0x%x", inputByte),
		"0", models.Setting.Contract)
	if err != nil {
		log.Append(fmt.Sprintf("调用合约失败 %s\n", err.Error()))
		return
	}

	if ret.Reverted {
		log.Append(fmt.Sprintf("合约执行终止 %s\n", ret.VMErr))
		return
	}

	// 对返回的结果进行解包
	// 去除0x
	retData, err := hex.DecodeString(ret.Data[2:])
	if err != nil {
		log.Append(fmt.Sprintf("对结果进行十六进制转换失败 %s\n", err.Error()))
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
			retI = append(argsI, &addr)

		default:
			var tmp string
			retI = append(argsI, &tmp)
		}
	}

	if len(retI) == 1 {
		ret := retI[0]
		err = abiObj.Unpack(ret, method, retData)
		if err != nil {
			log.Append(fmt.Sprintf("对调用结果进行解码失败 %s\n", err.Error()))
			return
		}
		arg := abiObj.Methods[method].Outputs[0]
		data, _ := json.Marshal(ret)
		log.Append(fmt.Sprintf("参数[%d](%s): %v \n", 0, arg.Name, string(data)))
		return
	}

	err = abiObj.Unpack(&retI, method, retData)
	if err != nil {
		log.Append(fmt.Sprintf("对调用结果进行解码失败 %s\n", err.Error()))
		return
	}

	for i, arg := range abiObj.Methods[method].Outputs {
		log.Append(fmt.Sprintf("参数[%d](%s): %v \n", i, arg.Name, retI[i]))
	}
	return
}