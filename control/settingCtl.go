package control

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/wupeaking/vechaintool/models"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 保存配置信息
func SaveSetting(privKey, contract, rpc, abiFile string, mainWin *ui.Window) {
	// 加载私钥
	if len(privKey) == 66 {
		privKey = privKey[2:]
	} else if len(privKey) != 64 {
		ui.MsgBoxError(mainWin, "错误提示", "Private key format error")
		return
	} else {
		_, err := hex.DecodeString(privKey)
		if err != nil {
			ui.MsgBoxError(mainWin, "错误提示", "Private key decoding failed")
			return
		}
	}
	// 检查合约地址是否正确
	if len(contract) != 42 {
		ui.MsgBoxError(mainWin, "错误提示", "The contract address is in the wrong format")
		return
	}

	if !strings.HasPrefix(rpc, "http://") {
		ui.MsgBoxError(mainWin, "错误提示", "Server address format error, "+rpc)
		return
	}

	var bufs []byte
	if abiFile != "" {
		// 尝试打开abi文件
		file, err := os.Open(abiFile)
		if err != nil {
			ui.MsgBoxError(mainWin, "错误提示", "Failed to open ABI file")
			return
		}
		_, err = abi.JSON(file)
		if err != nil {
			ui.MsgBoxError(mainWin, "错误提示", "Abi file deserialization failed")
			return
		}
		file.Seek(0, io.SeekStart)
		bufs, _ = ioutil.ReadAll(file)
	}
	models.Setting.PrivKey = privKey
	models.Setting.Contract = contract
	models.Setting.RPC = rpc
	models.Setting.ABI = string(bufs)

	content, _ := json.Marshal(models.Setting)
	settingFile, err := os.Create("./.vechain_setting.json")
	if err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "saving configuration file failed")
		return
	}
	_, err = settingFile.Write(content)
	if err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "saving configuration file failed")
		return
	}

	if err := models.Setting.Load(); err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "saving configuration file failed: "+err.Error())
		return
	}

	ui.MsgBox(mainWin, "提示", "saving configuration file success")
}

//LoadSetting 加载配置信息
func LoadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry *ui.Entry, mainWin *ui.Window) {

	// 尝试打开配置文件
	file, err := os.Open("./.vechain_setting.json")
	if err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "opening configuration file failed")
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "reading configuration file failed")
		return
	}

	if json.Unmarshal(content, &models.Setting) != nil {
		ui.MsgBoxError(mainWin, "错误提示", "Deserialization configuration file failed")
		return
	}

	privKeyEntry.SetText(models.Setting.PrivKey)
	contractEntry.SetText(models.Setting.Contract)
	rpcEntry.SetText(models.Setting.RPC)

	if err := models.Setting.Load(); err != nil {
		ui.MsgBoxError(mainWin, "提示", "loading configuration file failed: "+err.Error())
		return
	}

	ui.MsgBox(mainWin, "提示", "loading configuration file success")
}

//TryLoadSetting 尝试加载配置
func TryLoadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry *ui.Entry, mainWin *ui.Window) {

	// 尝试打开配置文件
	file, err := os.Open("./.vechain_setting.json")
	if err != nil {
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	if json.Unmarshal(content, &models.Setting) != nil {
		return
	}

	privKeyEntry.SetText(models.Setting.PrivKey)
	contractEntry.SetText(models.Setting.Contract)
	rpcEntry.SetText(models.Setting.RPC)
	if err := models.Setting.Load(); err != nil {
		fmt.Println("load err: ", err.Error())
	}
}
