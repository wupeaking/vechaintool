package models

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// 配置信息
type Configure struct {
	PrivKey        string `json:"priv_key"`
	Contract       string `json:"contract"`
	RPC            string `json:"rpc"`
	ABI            string `json:"abi"`
	abiObj         *abi.ABI
	priv           *ecdsa.PrivateKey
	Addr           common.Address
	ViewFuncDesc   []FuncDesc
	StatusFuncDesc []FuncDesc
}

var Setting Configure

func (cfg *Configure) Load() error {
	// 清空
	cfg.abiObj = nil
	cfg.priv = nil
	cfg.Addr = common.Address{}
	cfg.ViewFuncDesc = nil
	cfg.StatusFuncDesc = nil

	if err := cfg.loadPrivKey(); err != nil {
		return err
	}
	if err := cfg.loadABI(); err != nil {
		return err
	}
	cfg.viewsFunc()
	cfg.statusFunc()

	return nil
}

func (cfg *Configure) loadPrivKey() error {
	p, err := crypto.HexToECDSA(cfg.PrivKey)
	cfg.priv = p
	cfg.Addr = crypto.PubkeyToAddress(p.PublicKey)

	return err
}

func (cfg *Configure) loadABI() error {
	if cfg.ABI != "" {
		buf := bytes.NewBufferString(cfg.ABI)
		a, err := abi.JSON(buf)
		if err != nil {
			return err
		}
		cfg.abiObj = &a
		return nil
	}
	return nil
}

func (cfg *Configure) PrivateKey() *ecdsa.PrivateKey {
	return cfg.priv
}

func (cfg *Configure) ABIObj() *abi.ABI {
	return cfg.abiObj
}

type FuncDesc struct {
	Name   string
	Inputs []InputsDesc
}
type InputsDesc struct {
	ArgName string
	ArgType string
}

func (cfg *Configure) viewsFunc() []FuncDesc {
	if cfg.abiObj == nil {
		return nil
	}
	tmp := make([]FuncDesc, 0)

	for name, method := range cfg.abiObj.Methods {
		if !method.Const {
			continue
		}
		desc := FuncDesc{}
		desc.Name = name
		desc.Inputs = make([]InputsDesc, 0, len(method.Inputs))
		for i := range method.Inputs {
			in := InputsDesc{}
			in.ArgType = method.Inputs[i].Type.String()
			in.ArgName = method.Inputs[i].Name
			desc.Inputs = append(desc.Inputs, in)
		}
		tmp = append(tmp, desc)
	}
	cfg.ViewFuncDesc = tmp
	return tmp
}

func (cfg *Configure) statusFunc() []FuncDesc {
	if cfg.abiObj == nil {
		return nil
	}
	tmp := make([]FuncDesc, 0)

	for name, method := range cfg.abiObj.Methods {
		if method.Const {
			continue
		}
		desc := FuncDesc{}
		desc.Name = name
		desc.Inputs = make([]InputsDesc, 0, len(method.Inputs))
		for i := range method.Inputs {
			in := InputsDesc{}
			in.ArgType = method.Inputs[i].Type.String()
			in.ArgName = method.Inputs[i].Name
			desc.Inputs = append(desc.Inputs, in)
		}
		tmp = append(tmp, desc)
	}
	cfg.StatusFuncDesc = tmp
	return tmp
}
