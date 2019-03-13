package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/ethereum/go-ethereum/accounts/abi"
	// "github.com/ethereum/go-ethereum/crypto"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var mainwin *ui.Window

var setting = struct {
	PrivKey  string `json:"priv_key"`
	Contract string `json:"contract"`
	RPC      string `json:"rpc"`
	ABI      string `json:"abi"`
}{}

func makeTransferPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	// 添加发起方
	//senderhbox := ui.NewHorizontalBox()
	//senderhbox.SetPadded(false)
	//vbox.Append(senderhbox, false)
	//
	//_, err := crypto.HexToECDSA(setting.PrivKey)
	//if err != nil {
	//	msg.SetText("解析配置中的私钥出错"+err.Error())
	//}
	//
	//sender := ui.NewEntry()
	//
	//sender.SetText("0xd2a58c88c60593b5af1ec177b7cea838d261fd0d")
	//sender.SetReadOnly(true)
	//senderhbox.Append(ui.NewLabel("sender:   "), false)
	//senderhbox.Append(sender, true)


	// 添加接收方
	tohbox := ui.NewHorizontalBox()
	tohbox.SetPadded(true)
	vbox.Append(tohbox, false)
	to := ui.NewEntry()
	tohbox.Append(ui.NewLabel("to:           "), false)
	tohbox.Append(to, true)

	// 添加金额
	amounthbox := ui.NewHorizontalBox()
	amounthbox.SetPadded(true)
	vbox.Append(amounthbox, false)
	amount := ui.NewEntry()
	amounthbox.Append(ui.NewLabel("amount:  "), false)
	amounthbox.Append(amount, true)

	// 选择币种类型
	currecyhbox := ui.NewHorizontalBox()
	currecyhbox.SetPadded(true)
	vbox.Append(currecyhbox, false)
	cbox := ui.NewCombobox()
	cbox.Append("VET")
	cbox.Append("VTHO")
	cbox.Append("ERC20")
	currecyhbox.Append(ui.NewLabel("currency type:"), false)
	currecyhbox.Append(cbox, false)

	// 选择币种类型
	erc20Addr := ui.NewEntry()
	currecyhbox.Append(ui.NewLabel("token address:"), false)
	currecyhbox.Append(erc20Addr, true)

	// 设置费用
	spinbox := ui.NewSpinbox(0, 255)
	slider := ui.NewSlider(0, 255)

	spinbox.OnChanged(func(*ui.Spinbox) {
		slider.SetValue(spinbox.Value())
	})
	slider.OnChanged(func(*ui.Slider) {
		spinbox.SetValue(slider.Value())
	})
	gashbox := ui.NewHorizontalBox()
	gashbox.SetPadded(true)
	vbox.Append(gashbox, false)
	gashbox.Append(ui.NewLabel("fee:        "), false)
	gashbox.Append(slider, false)
	gashbox.Append(spinbox, false)


	txBtn := ui.NewButton("转账")
	txBtn.OnClicked(func(*ui.Button) {

	})
	vbox.Append(txBtn, true)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	//
	oplog := ui.NewMultilineEntry()
	oplog.SetReadOnly(true)
	vbox.Append(oplog, true)
	return vbox
}


// 合约函数调用页
func makeContractPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	group := ui.NewGroup("view function")
	group.SetMargined(true)
	hbox.Append(group, true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	viewbox := ui.NewCombobox()
	viewbox.Append("VET")
	viewbox.Append("VTHO")
	viewbox.Append("ERC20")
	vbox.Append(viewbox, false)

	// todo:: 添加参数列表
	argsGroup := ui.NewGroup("args list")
	vbox.Append(argsGroup, false)


	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	privKeyEntry := ui.NewPasswordEntry()
	contractEntry := ui.NewEntry()
	rpcEntry := ui.NewEntry()
	entryForm.Append("私钥(十六进制格式)", privKeyEntry, false)
	entryForm.Append("智能合约地址", contractEntry, false)
	entryForm.Append("唯链客户端地址", rpcEntry, false)
	argsGroup.SetChild(entryForm)
	callBtn := ui.NewButton("call")
	callResult := ui.NewMultilineEntry()
	callResult.Append("1\n")
	callResult.Append("2")

	callBtn.OnClicked(func(button *ui.Button) {
		//todo:: 调用视图函数
	})
	vbox.Append(callBtn, false)

	vbox.Append(callResult, false)

	//------------- 状态函数

	statusGroup := ui.NewGroup("status function")
	statusGroup.SetMargined(true)
	hbox.Append(statusGroup, true)

	statusVbox := ui.NewVerticalBox()
	statusVbox.SetPadded(true)
	statusGroup.SetChild(statusVbox)

	statusviewbox := ui.NewCombobox()
	statusviewbox.Append("VET")
	statusviewbox.Append("VTHO")
	statusviewbox.Append("ERC20")
	statusVbox.Append(statusviewbox, false)

	// todo:: 添加参数列表
	statusArgsGroup := ui.NewGroup("args list")
	statusVbox.Append(statusArgsGroup, false)

	statusForm := ui.NewForm()
	statusForm.SetPadded(true)
	statusForm.Append("私钥(十六进制格式)", ui.NewEntry(), false)
	statusForm.Append("智能合约地址", ui.NewEntry(), false)
	statusForm.Append("唯链客户端地址", ui.NewEntry(), false)
	statusArgsGroup.SetChild(statusForm)


	txBtn := ui.NewButton("start transact")
	txResult := ui.NewMultilineEntry()
	txResult.Append("1\n")
	txResult.Append("2")

	txBtn.OnClicked(func(button *ui.Button) {
		//todo:: 调用视图函数
	})
	statusVbox.Append(txBtn, false)

	statusVbox.Append(txResult, false)


	return hbox
}

func makeDataChoosersPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	vbox.Append(ui.NewDatePicker(), false)
	vbox.Append(ui.NewTimePicker(), false)
	vbox.Append(ui.NewDateTimePicker(), false)
	vbox.Append(ui.NewFontButton(), false)
	vbox.Append(ui.NewColorButton(), false)

	hbox.Append(ui.NewVerticalSeparator(), false)

	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, true)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	button := ui.NewButton("Open File")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry.SetText(filename)
	})
	grid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry,
		1, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Save File")
	entry2 := ui.NewEntry()
	entry2.SetReadOnly(true)
	button.OnClicked(func(*ui.Button) {
		filename := ui.SaveFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry2.SetText(filename)
	})
	grid.Append(button,
		0, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry2,
		1, 1, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	msggrid := ui.NewGrid()
	msggrid.SetPadded(true)
	grid.Append(msggrid,
		0, 2, 2, 1,
		false, ui.AlignCenter, false, ui.AlignStart)

	button = ui.NewButton("Message Box")
	button.OnClicked(func(*ui.Button) {
		ui.MsgBox(mainwin,
			"This is a normal message box.",
			"More detailed information can be shown here.")
	})
	msggrid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	button = ui.NewButton("Error Box")
	button.OnClicked(func(*ui.Button) {
		ui.MsgBoxError(mainwin,
			"This message box describes an error.",
			"More detailed information can be shown here.")
	})
	msggrid.Append(button,
		1, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	return hbox
}

func makeSettingPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	group := ui.NewGroup("设置列表")
	group.SetMargined(true)
	vbox.Append(group, false)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	privKeyEntry := ui.NewPasswordEntry()
	contractEntry := ui.NewEntry()
	rpcEntry := ui.NewEntry()
	entryForm.Append("私钥(十六进制格式)", privKeyEntry, false)
	entryForm.Append("智能合约地址", contractEntry, false)
	entryForm.Append("唯链客户端地址", rpcEntry, false)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	entryForm.Append("", grid, false)

	abiBtn := ui.NewButton("打开ABI文件")
	abiEntry := ui.NewEntry()
	abiEntry.SetReadOnly(true)
	abiBtn.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "未选择"
		}
		abiEntry.SetText(filename)
	})
	grid.Append(abiBtn,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(abiEntry,
		1, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	btngrid := ui.NewGrid()
	// btngrid.SetPadded(true)
	vbox.Append(btngrid, false)

	saveBtn := ui.NewButton("保存配置")
	saveBtn.OnClicked(func(*ui.Button) {
		saveSetting(privKeyEntry.Text(), contractEntry.Text(), rpcEntry.Text(), abiEntry.Text())
	})
	btngrid.Append(ui.NewLabel("                                             "),
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	btngrid.Append(saveBtn,
		1, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	loadBtn := ui.NewButton("加载配置")
	loadBtn.OnClicked(func(*ui.Button) {
		loadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry)
	})
	btngrid.Append(loadBtn,
		2, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	return vbox
}

func saveSetting(privKey, contract, rpc, abiFile string) {
	setting.PrivKey = privKey
	setting.Contract = contract
	setting.RPC = rpc
	// 加载私钥
	if len(privKey) == 66 {
		privKey = privKey[2:]
	} else if len(privKey) != 64 {
		ui.MsgBoxError(mainwin, "错误提示", "私钥格式错误")
		return
	} else {
		_, err := hex.DecodeString(privKey)
		if err != nil {
			ui.MsgBoxError(mainwin, "错误提示", "私钥解码失败")
			return
		}
	}
	// 检查合约地址是否正确
	if len(contract) != 42 {
		ui.MsgBoxError(mainwin, "错误提示", "合约地址格式错误")
		return
	}

	if !strings.HasPrefix(rpc, "http://") {
		ui.MsgBoxError(mainwin, "错误提示", "服务器地址格式错误, "+rpc)
		return
	}

	// 尝试打开abi文件
	file, err := os.Open(abiFile)
	if err != nil {
		ui.MsgBoxError(mainwin, "错误提示", "打开ABI文件失败")
		return
	}

	_, err = abi.JSON(file)
	if err != nil {
		ui.MsgBoxError(mainwin, "错误提示", "abi文件反序列化失败")
		return
	}
	file.Seek(0, io.SeekStart)
	bytes, _ := ioutil.ReadAll(file)
	setting.ABI = string(bytes)
	content, _ := json.Marshal(setting)

	settingFile, err := os.Create("./.vechain_setting.json")
	if err != nil {
		ui.MsgBoxError(mainwin, "错误提示", "保存配置文件出错")
		return
	}
	_, err = settingFile.Write(content)
	if err != nil {
		ui.MsgBoxError(mainwin, "错误提示", "保存配置文件出错")
		return
	}
	ui.MsgBox(mainwin, "提示", "保存配置文件成功")
}

func loadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry *ui.Entry) {

	// 尝试打开配置文件
	file, err := os.Open("./.vechain_setting.json")
	if err != nil {
		ui.MsgBoxError(mainwin, "错误提示", "打开配置文件出错")
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		ui.MsgBoxError(mainwin, "错误提示", "读取配置文件出错")
		return
	}

	if json.Unmarshal(content, &setting) != nil {
		ui.MsgBoxError(mainwin, "错误提示", "反序列化配置文件出错")
		return
	}

	privKeyEntry.SetText(setting.PrivKey)
	contractEntry.SetText(setting.Contract)
	rpcEntry.SetText(setting.RPC)

	ui.MsgBox(mainwin, "提示", "加载配置文件成功")
}


func setupUI() {
	mainwin = ui.NewWindow("唯链调试小工具", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)


	tab.Append("转账", makeTransferPage())
	tab.SetMargined(0, true)

	tab.Append("合约", makeContractPage())
	tab.SetMargined(1, true)

	tab.Append("编码", makeDataChoosersPage())
	tab.SetMargined(2, true)

	tab.Append("设置", makeSettingPage())
	tab.SetMargined(3, true)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
