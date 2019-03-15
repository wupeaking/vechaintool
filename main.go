package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/wupeaking/vechaintool/control"
	"github.com/wupeaking/vechaintool/models"
	// "github.com/ethereum/go-ethereum/crypto"
)

var mainwin *ui.Window
var refresh = make(chan struct{}, 1)

func makeTransferPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

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

	currecy := ""
	cbox.OnSelected(func(combobox *ui.Combobox) {
		switch combobox.Selected() {
		case 0:
			currecy = "VET"
		case 1:
			currecy = "VTHO"
		case 2:
			currecy = "ERC20"
		}
	})

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
	gashbox.Append(ui.NewLabel("GasPriceCoef:    "), false)
	gashbox.Append(slider, false)
	gashbox.Append(spinbox, false)

	txBtn := ui.NewButton("转账")
	vbox.Append(txBtn, true)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	oplog := ui.NewMultilineEntry()
	oplog.SetReadOnly(true)
	vbox.Append(oplog, true)

	txBtn.OnClicked(func(*ui.Button) {
		control.Transfer(currecy, amount.Text(), to.Text(), erc20Addr.Text(), oplog)
	})
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
	for _, fd := range models.Setting.ViewFuncDesc {
		viewbox.Append(fd.Name)
	}
	vbox.Append(viewbox, false)

	argsGroup := ui.NewGroup("args list")
	vbox.Append(argsGroup, false)

	callEntrys := make([]*ui.Entry, 0)
	curViewFuncName := ""
	viewbox.OnSelected(func(combobox *ui.Combobox) {
		form := ui.NewForm()
		form.SetPadded(true)
		callEntrys = callEntrys[:0]
		funcIndex := combobox.Selected()
		if funcIndex >= len(models.Setting.ViewFuncDesc) || funcIndex < 0 {
			return
		}
		curViewFuncName = models.Setting.ViewFuncDesc[funcIndex].Name

		for _, input := range models.Setting.ViewFuncDesc[funcIndex].Inputs {
			e := ui.NewEntry()
			form.Append(fmt.Sprintf("%s(%s)", input.ArgName, input.ArgType), e, false)
			callEntrys = append(callEntrys, e)
		}
		argsGroup.SetChild(form)
	})

	callBtn := ui.NewButton("call")
	callResult := ui.NewMultilineEntry()

	callBtn.OnClicked(func(button *ui.Button) {
		callResult.SetText("")
		args := make([]string, 0, len(callEntrys))
		for i := range callEntrys {
			args = append(args, callEntrys[i].Text())
		}
		control.ViewFuncCall(curViewFuncName, args, callResult)
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
	for _, fd := range models.Setting.StatusFuncDesc {
		statusviewbox.Append(fd.Name)
	}
	statusVbox.Append(statusviewbox, false)

	statusArgsGroup := ui.NewGroup("args list")
	statusVbox.Append(statusArgsGroup, false)

	entrys := make([]*ui.Entry, 0)
	curFuncName := ""
	statusviewbox.OnSelected(func(combobox *ui.Combobox) {
		form := ui.NewForm()
		form.SetPadded(true)
		entrys = entrys[:0]
		funcIndex := combobox.Selected()
		if funcIndex >= len(models.Setting.StatusFuncDesc) || funcIndex < 0 {
			return
		}
		curFuncName = models.Setting.StatusFuncDesc[funcIndex].Name

		for _, input := range models.Setting.StatusFuncDesc[funcIndex].Inputs {
			e := ui.NewEntry()
			form.Append(fmt.Sprintf("%s(%s)", input.ArgName, input.ArgType), e, false)
			entrys = append(entrys, e)
		}
		statusArgsGroup.SetChild(form)
	})

	txBtn := ui.NewButton("start transact")
	txResult := ui.NewMultilineEntry()

	txBtn.OnClicked(func(button *ui.Button) {
		txResult.SetText("")
		args := make([]string, 0, len(entrys))
		for i := range entrys {
			args = append(args, entrys[i].Text())
		}
		control.CallContract(curFuncName, args, txResult)
	})
	statusVbox.Append(txBtn, false)
	statusVbox.Append(txResult, false)
	return hbox
}

func makeEncodingPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)
	log := ui.NewMultilineEntry()
	log.SetReadOnly(true)

	// ABI 编码字符串
	strEntry := ui.NewMultilineEntry()

	abiEncodeBtn := ui.NewButton("ABI encode")
	abiEncodeBtn.OnClicked(func(*ui.Button) {
		control.ABIPack(strEntry.Text(), log)
	})

	abiDecodeBtn := ui.NewButton("ABI decode")
	abiDecodeBtn.OnClicked(func(*ui.Button) {
		control.ABIUnPack(strEntry.Text(), log)
	})

	signBtn := ui.NewButton("signature")
	signBtn.OnClicked(func(*ui.Button) {
	})

	hash256Btn := ui.NewButton("SHA256")
	hash256Btn.OnClicked(func(*ui.Button) {
	})

	keccak256Btn := ui.NewButton("Keccak256")
	keccak256Btn.OnClicked(func(*ui.Button) {
	})

	base64Encode := ui.NewButton("b64 encode")
	base64Encode.OnClicked(func(*ui.Button) {
		control.Base64Encode(strEntry.Text(), log)
	})

	base64Decode := ui.NewButton("b64 decode")
	base64Decode.OnClicked(func(*ui.Button) {
		control.Base64Decode(strEntry.Text(), log)
	})

	grid.Append(strEntry,
		0, 0, 2, 7,
		true, ui.AlignFill, true, ui.AlignFill)
	grid.Append(abiEncodeBtn,
		1, 0, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(abiDecodeBtn,
		1, 1, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(signBtn,
		1, 2, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(hash256Btn,
		1, 3, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(keccak256Btn,
		1, 4, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(base64Encode,
		1, 5, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(base64Decode,
		1, 6, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	vbox.Append(ui.NewHorizontalSeparator(), false)


	vbox.Append(log,  true)

	return vbox
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
	btngrid.SetPadded(true)
	vbox.Append(btngrid, true)

	msgGrid := ui.NewGrid()
	msgGrid.SetPadded(true)
	btngrid.Append(msgGrid,
		0, 0, 2, 1,
		false, ui.AlignCenter, false, ui.AlignStart)

	saveBtn := ui.NewButton("保存配置")
	saveBtn.OnClicked(func(*ui.Button) {
		control.SaveSetting(privKeyEntry.Text(), contractEntry.Text(), rpcEntry.Text(), abiEntry.Text(), mainwin)
		refresh <- struct{}{}
	})

	msgGrid.Append(saveBtn,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	loadBtn := ui.NewButton("加载配置")
	loadBtn.OnClicked(func(*ui.Button) {
		control.LoadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry, mainwin)
		refresh <- struct{}{}
	})
	msgGrid.Append(loadBtn,
		1, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	control.TryLoadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry, mainwin)
	return vbox
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

	tab.Append("设置", makeSettingPage())
	tab.SetMargined(0, true)

	tab.Append("转账", makeTransferPage())
	tab.SetMargined(1, true)

	tab.Append("合约", makeContractPage())
	tab.SetMargined(2, true)

	tab.Append("编码", makeEncodingPage())
	tab.SetMargined(3, true)

	go func() {
		for {
			select {
			case <-refresh:
				ui.QueueMain(func() {
					tab.Delete(2)
					tab.InsertAt("合约", 2, makeContractPage())
				})
			}
		}
	}()
	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
