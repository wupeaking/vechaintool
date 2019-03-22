package view

import (
	"fmt"
	"github.com/andlabs/ui"
	"github.com/wupeaking/vechaintool/control"
	"github.com/wupeaking/vechaintool/models"
)

//MakeContractPage 合约函数调用页
func MakeContractPage(window *ui.Window) ui.Control {
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
