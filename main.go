package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	hbox.Append(ui.NewButton("Button"), false)
	hbox.Append(ui.NewCheckbox("Checkbox"), false)

	vbox.Append(ui.NewLabel("This is a label. Right now, labels can only span one line."), false)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	group := ui.NewGroup("Entries")
	group.SetMargined(true)
	vbox.Append(group, true)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	entryForm.Append("Entry", ui.NewEntry(), false)
	entryForm.Append("Password Entry", ui.NewPasswordEntry(), false)
	entryForm.Append("Search Entry", ui.NewSearchEntry(), false)
	entryForm.Append("Multiline Entry", ui.NewMultilineEntry(), true)
	entryForm.Append("Multiline Entry No Wrap", ui.NewNonWrappingMultilineEntry(), true)

	return vbox
}

func makeNumbersPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	group := ui.NewGroup("Numbers")
	group.SetMargined(true)
	hbox.Append(group, true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	spinbox := ui.NewSpinbox(0, 100)
	slider := ui.NewSlider(0, 100)
	pbar := ui.NewProgressBar()
	spinbox.OnChanged(func(*ui.Spinbox) {
		slider.SetValue(spinbox.Value())
		pbar.SetValue(spinbox.Value())
	})
	slider.OnChanged(func(*ui.Slider) {
		spinbox.SetValue(slider.Value())
		pbar.SetValue(slider.Value())
	})
	vbox.Append(spinbox, false)
	vbox.Append(slider, false)
	vbox.Append(pbar, false)

	ip := ui.NewProgressBar()
	ip.SetValue(-1)
	vbox.Append(ip, false)

	group = ui.NewGroup("Lists")
	group.SetMargined(true)
	hbox.Append(group, true)

	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	cbox := ui.NewCombobox()
	cbox.Append("Combobox Item 1")
	cbox.Append("Combobox Item 2")
	cbox.Append("Combobox Item 3")
	vbox.Append(cbox, false)

	ecbox := ui.NewEditableCombobox()
	ecbox.Append("Editable Item 1")
	ecbox.Append("Editable Item 2")
	ecbox.Append("Editable Item 3")
	vbox.Append(ecbox, false)

	rb := ui.NewRadioButtons()
	rb.Append("Radio Button 1")
	rb.Append("Radio Button 2")
	rb.Append("Radio Button 3")
	vbox.Append(rb, false)

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

	tab.Append("转账", makeBasicControlsPage())
	tab.SetMargined(0, true)

	tab.Append("合约", makeNumbersPage())
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
