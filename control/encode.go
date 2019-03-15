package control

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/wupeaking/vechaintool/models"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
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
	log.SetText("")
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
	log.SetText("")
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

func Signature(content string, log *ui.MultilineEntry) {
	log.SetText("")
	if models.Setting.PrivateKey() == nil {
		log.Append(fmt.Sprintf("signature failed: %s", "private key is unknown"))
		return
	}
	d := sha3.NewLegacyKeccak256()
	d.Write([]byte(content))
	hash := d.Sum(nil)
	sign, err := crypto.Sign(hash, models.Setting.PrivateKey())
	if err != nil {
		log.Append(fmt.Sprintf("signature failed: %s", err.Error()))
		return
	}

	log.Append(fmt.Sprintf("signature result(hex): %0x", sign))
	return
}

func Keccak256Hash(content string, log *ui.MultilineEntry) {
	log.SetText("")
	d := sha3.NewLegacyKeccak256()
	d.Write([]byte(content))
	hash := d.Sum(nil)

	log.Append(fmt.Sprintf("Keccak256 hash result(hex): %0x", hash))
	return
}

func Sha256(content string, log *ui.MultilineEntry) {
	log.SetText("")
	d := sha256.New()
	d.Write([]byte(content))
	hash := d.Sum(nil)
	log.Append(fmt.Sprintf("sha256 hash result(hex): %0x", hash))
	return
}