package view

import (
	"github.com/andlabs/ui"
	"github.com/wupeaking/vechaintool/control"
)

//MakeSettingPage 生成配置页
func MakeSettingPage(mainwin *ui.Window) ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	group := ui.NewGroup("设置页列表")
	group.SetMargined(true)
	vbox.Append(group, false)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	privKeyEntry := ui.NewPasswordEntry()
	contractEntry := ui.NewEntry()
	rpcEntry := ui.NewEntry()
	entryForm.Append("私钥(十六进制)", privKeyEntry, false)
	entryForm.Append("合约地址", contractEntry, false)
	entryForm.Append("RPC服务地址", rpcEntry, false)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	entryForm.Append("", grid, false)

	abiBtn := ui.NewButton("打开ABI文件")
	abiEntry := ui.NewEntry()
	abiEntry.SetReadOnly(true)
	abiBtn.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "unselected"
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

	saveBtn := ui.NewButton("保存设置")
	saveBtn.OnClicked(func(*ui.Button) {
		control.SaveSetting(privKeyEntry.Text(), contractEntry.Text(), rpcEntry.Text(), abiEntry.Text(), mainwin)
		RefreshfPage("settingPage")
	})

	msgGrid.Append(saveBtn,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	loadBtn := ui.NewButton("加载设置")
	loadBtn.OnClicked(func(*ui.Button) {
		control.LoadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry, mainwin)
		RefreshfPage("settingPage")
	})
	msgGrid.Append(loadBtn,
		1, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	control.TryLoadSetting(privKeyEntry, contractEntry, rpcEntry, abiEntry, mainwin)
	return vbox
}
