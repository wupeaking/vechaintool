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
		ui.MsgBoxError(mainWin, "错误提示", "私钥格式错误")
		return
	} else {
		_, err := hex.DecodeString(privKey)
		if err != nil {
			ui.MsgBoxError(mainWin, "错误提示", "私钥解码失败")
			return
		}
	}
	// 检查合约地址是否正确
	if len(contract) != 42 {
		ui.MsgBoxError(mainWin, "错误提示", "合约地址格式错误")
		return
	}

	if !strings.HasPrefix(rpc, "http://") {
		ui.MsgBoxError(mainWin, "错误提示", "服务器地址格式错误, "+rpc)
		return
	}

	var bufs []byte
	if abiFile != "" {
		// 尝试打开abi文件
		file, err := os.Open(abiFile)
		if err != nil {
			ui.MsgBoxError(mainWin, "错误提示", "打开ABI文件失败")
			return
		}
		_, err = abi.JSON(file)
		if err != nil {
			ui.MsgBoxError(mainWin, "错误提示", "abi文件反序列化失败")
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
		ui.MsgBoxError(mainWin, "错误提示", "保存配置文件出错")
		return
	}
	_, err = settingFile.Write(content)
	if err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "保存配置文件出错")
		return
	}

	if err := models.Setting.Load(); err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "保存配置出错: "+err.Error())
		return
	}

	ui.MsgBox(mainWin, "提示", "保存配置文件成功")
}

//LoadSetting 加载配置信息
func LoadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry *ui.Entry, mainWin *ui.Window) {

	// 尝试打开配置文件
	file, err := os.Open("./.vechain_setting.json")
	if err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "打开配置文件出错")
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		ui.MsgBoxError(mainWin, "错误提示", "读取配置文件出错")
		return
	}

	if json.Unmarshal(content, &models.Setting) != nil {
		ui.MsgBoxError(mainWin, "错误提示", "反序列化配置文件出错")
		return
	}

	privKeyEntry.SetText(models.Setting.PrivKey)
	contractEntry.SetText(models.Setting.Contract)
	rpcEntry.SetText(models.Setting.RPC)

	if err := models.Setting.Load(); err != nil {
		ui.MsgBoxError(mainWin, "提示", "加载配置出错: "+err.Error())
		return
	}

	ui.MsgBox(mainWin, "提示", "加载配置文件成功")
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
