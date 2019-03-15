package control

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
	"github.com/wupeaking/vechaintool/models"
	"github.com/wupeaking/vechaintool/vechainclient"
	"github.com/wupeaking/vechaintool/view"
	"math/big"
	"strings"
	"time"
)

func Transfer(currecy, amount, to, contract string, log *ui.MultilineEntry) {
	if currecy == "" {
		log.Append("The correct currency type was not selected\n")
		return
	}

	value, ok := new(big.Int).SetString(amount, 0)
	if !ok {
		log.Append("amount format exception\n")
		return
	}

	if models.Setting.Addr.String() == "" {
		log.Append("Private key not set\n")
		return
	}

	if len(to) != 42 {
		log.Append("Receiver address error\n")
		return
	}

	if strings.ToLower(currecy) == "erc20" && len(contract) != 42 {
		log.Append("not fill in the contract address\n")
		return
	}

	// 查询当前区块号
	veCli, _ := vechain.NewVeChainClient(models.Setting.RPC, "", "", 10)

	blk, err := veCli.BlockInfo(context.Background(), 0)
	if err != nil {
		return
	}

	gasUsed, err := calcTransferGasUsed(currecy, value, veCli, contract)
	if err != nil {
		log.Append(fmt.Sprintf("Analog calculation gas error, err: %s\n", err.Error()))
		return
	}

	allVtho, err := queryBalance("vtho", models.Setting.Addr.String(), contract, veCli)
	if err != nil {
		log.Append(fmt.Sprintf("querying the balance of vtho failed, err: %s\n", err.Error()))
		return
	}
	// gas --> vtho(wei)  (gas/1000)*10^18=gas*10^15
	totalUsedVtho := new(big.Int).Mul(big.NewInt(int64(gasUsed)), big.NewInt(10e15))
	if allVtho.Cmp(totalUsedVtho) < 0 {
		log.Append(fmt.Sprintf("Vtho is not enough to pay for miners\n"))
		return
	}

	// 构造交易 创建原始交易
	_, rawTx, err := constructTxData(models.Setting.PrivateKey(), currecy, uint32(blk.BlockNum),
		contract, to, gasUsed, value)

	if err != nil {
		log.Append(fmt.Sprintf("Constructing the original transaction error: %s\n", err.Error()))
		return
	}

	msg := fmt.Sprintf("Confirm sending  %s %s to %s, This transaction will cost Gas: %d",
		value.String(), currecy, to, gasUsed)
	view.ConfirmDialog(msg, func() {
		// 广播交易
		txid, err := veCli.PushTx(context.Background(), rawTx)
		if err != nil {
			log.Append(fmt.Sprintf("Broadcast transaction failed: %s\n", err.Error()))
			return
		}
		log.Append(fmt.Sprintf("Broadcast transaction success, tx_id: %s\n", txid))
		return
	}, func() {
		log.Append(fmt.Sprintf("cancel transaction\n"))
		return
	})
	return
}

//calcTransferGasUsed 模拟计算实际转账需要的gas
func calcTransferGasUsed(curreny string, amount *big.Int, cli *vechain.Client, erc20Addr string) (uint64, error) {
	// https://github.com/vechain/thor/wiki/FAQ#what-is-intrinsic-gas-
	if strings.ToLower(curreny) == "vet" {
		return 5000 + 16000, nil
	}

	//address _to, uint256 _value
	method := abi.Method{Name: "transfer", Const: false}
	addressType, _ := abi.NewType("address", nil)
	uin256Type, _ := abi.NewType("uint256", nil)

	_to := abi.Argument{Name: "_to", Type: addressType, Indexed: false}
	_value := abi.Argument{Name: "_value", Type: uin256Type, Indexed: false}
	method.Inputs = abi.Arguments{_to, _value}

	toAddr := common.HexToAddress("0xb9b7e0cb2edf5ea031c8b297a5a1fa20379b6a0a")
	argsData, _ := method.Inputs.Pack(toAddr, amount)
	data := append(method.Id(), argsData...)

	result, err := cli.SimulateContract(context.Background(), fmt.Sprintf("0x%0x", data),
		"0", erc20Addr)

	if err != nil {
		return 0, err
	}
	if result.Reverted {
		return 0, fmt.Errorf("虚拟机执行reverted")
	}
	// txGas + (clauses.type + dataGas + vmGas)*len(clauses)

	dataGas := func(input []byte) uint64 {
		const nzgas = 68
		return uint64(len(input)) * nzgas
	}

	return (result.GasUsed+16000+dataGas(data))*uint64(1) + 5000, nil
}

//queryBalance 查询本币或者token余额
func queryBalance(currency string, account string, erc20Addr string, cli *vechain.Client) (*big.Int, error) {
	currency = strings.ToLower(currency)
	if currency == "vet" || currency == "vtho" {
		balance, err := cli.BalanceByAddress(context.Background(), account)
		if err != nil {
			return nil, err
		}
		if currency == "vet" {
			return balance.Balance, nil
		}
		return balance.Energy, nil
	}

	method := abi.Method{Name: "balanceOf", Const: false}
	addressType, _ := abi.NewType("address", nil)
	uin256Type, _ := abi.NewType("uint256", nil)

	_owner := abi.Argument{Name: "_owner", Type: addressType, Indexed: false}
	_value := abi.Argument{Name: "balance", Type: uin256Type, Indexed: false}
	method.Inputs = abi.Arguments{_owner}
	method.Outputs = abi.Arguments{_value}

	toAddr := common.HexToAddress(account)
	argsData, _ := method.Inputs.Pack(toAddr)
	input := append(method.Id(), argsData...)

	result, err := cli.SimulateContract(context.Background(), fmt.Sprintf("0x%0x", input),
		"0", erc20Addr)

	if err != nil {
		return nil, err
	}
	//value := reflect.New(reflect.TypeOf(big.NewInt(0)))
	value := big.NewInt(0)
	resultData, _ := hex.DecodeString(result.Data[2:])

	err = method.Outputs.Unpack(&value, resultData)
	if err != nil {
		return nil, err
	}
	// value.Elem().Interface().(*big.Int)
	return value, nil
}

// constructRawTransfer 构造交易 返回值： 交易ID 原始交易内容 错误
func constructTxData(prvk *ecdsa.PrivateKey, currency string, blockNum uint32,
	erc20Addr, to string, gas uint64, amount *big.Int) (string, string, error) {
	currency = strings.ToLower(currency)
	//   chaintag  创世区块ID 最后一个字节 测试链为0x27 生产链为0x4a
	trx := new(tx.Builder).ChainTag(0x27).
		BlockRef(tx.NewBlockRef(blockNum)).
		Expiration(720).
		GasPriceCoef(0).
		Gas(gas).
		DependsOn(nil).
		Nonce(uint64(time.Now().UnixNano()))

	var tokenAddr thor.Address
	if currency != "vet" {
		t, err := thor.ParseAddress(erc20Addr)
		if err != nil {
			return "", "", err
		}
		tokenAddr = t
	}

	toAddr, err := thor.ParseAddress(to)
	if err != nil {
		return "", "", err
	}
	if currency == "vet" {
		trx.Clause(tx.NewClause(&toAddr).WithValue(amount).WithData(nil))
	} else {
		// token 转账
		//address _to, uint256 _value
		method := abi.Method{Name: "transfer", Const: false}
		addressType, _ := abi.NewType("address", nil)
		uin256Type, _ := abi.NewType("uint256", nil)
		_to := abi.Argument{Name: "_to", Type: addressType, Indexed: false}
		_value := abi.Argument{Name: "_value", Type: uin256Type, Indexed: false}
		method.Inputs = abi.Arguments{_to, _value}
		argsData, _ := method.Inputs.Pack(toAddr, amount)
		input := append(method.Id(), argsData...)
		trx.Clause(tx.NewClause(&tokenAddr).WithValue(big.NewInt(0)).WithData(input))
	}

	trxBuild := trx.Build()
	sig, err := crypto.Sign(trxBuild.SigningHash().Bytes(), prvk)
	if err != nil {
		return "", "", err
	}

	trxBuild = trxBuild.WithSignature(sig)
	d, err := rlp.EncodeToBytes(trxBuild)
	if err != nil {
		return "", "", err
	}
	return trxBuild.ID().String(), hex.EncodeToString(d), nil
}
