package view

import (
	"github.com/andlabs/ui"
	"github.com/wupeaking/vechaintool/control"
)

//MakeSettingPage 生成配置页
func MakeSettingPage(mainwin *ui.Window) ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	group := ui.NewGroup("setting list")
	group.SetMargined(true)
	vbox.Append(group, false)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	privKeyEntry := ui.NewPasswordEntry()
	contractEntry := ui.NewEntry()
	rpcEntry := ui.NewEntry()
	entryForm.Append("private key(hex)", privKeyEntry, false)
	entryForm.Append("contract address", contractEntry, false)
	entryForm.Append("vechain rpc", rpcEntry, false)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	entryForm.Append("", grid, false)

	abiBtn := ui.NewButton("open abi file")
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

	saveBtn := ui.NewButton("save setting")
	saveBtn.OnClicked(func(*ui.Button) {
		control.SaveSetting(privKeyEntry.Text(), contractEntry.Text(), rpcEntry.Text(), abiEntry.Text(), mainwin)
		RefreshfPage("settingPage")
	})

	msgGrid.Append(saveBtn,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	loadBtn := ui.NewButton("load setting")
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
