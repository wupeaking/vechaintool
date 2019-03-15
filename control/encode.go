package control

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func Base64Encode(content string, log *ui.MultilineEntry) {
	log.SetText("")
	result := base64.StdEncoding.EncodeToString([]byte(content))
	log.Append(fmt.Sprintf("base64 encode result: %s", result))
}

func Base64Decode(content string, log *ui.MultilineEntry) {
	log.SetText("")
	result, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		log.Append(fmt.Sprintf("base64 decode failed: %s", err.Error()))
		return
	}
	log.Append(fmt.Sprintf("base64 decode result: %s", string(result)))
	return
}

func ABIPack(content string, log *ui.MultilineEntry) {
	strType, _ := abi.NewType("string", nil)
	abiArgs := abi.Arguments{abi.Argument{Name: "_", Type: strType, Indexed: false}}
	buf, err := abiArgs.Pack(content)
	if err != nil {
		log.Append(fmt.Sprintf("abi encode failed: %s", err.Error()))
		return
	}
	log.Append(fmt.Sprintf("abi encode result(hex): %0x", buf))
	return
}

func ABIUnPack(content string, log *ui.MultilineEntry) {
	data, err := hex.DecodeString(content)
	if err != nil {
		log.Append(fmt.Sprintf("input must be hex string: %s", err.Error()))
		return
	}

	strType, _ := abi.NewType("string", nil)
	abiArgs := abi.Arguments{abi.Argument{Name: "_", Type: strType, Indexed: false}}
	var result string
	err = abiArgs.Unpack(&result, data)
	if err != nil {
		log.Append(fmt.Sprintf("abi decode failed: %s", err.Error()))
		return
	}
	log.Append(fmt.Sprintf("abi encode result: %s", result))
	return
}