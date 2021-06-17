package view

import (
	"github.com/andlabs/ui"
	"github.com/wupeaking/vechaintool/control"
)

//MakeTransferPage 生成转账页
func MakeTransferPage(mainwin *ui.Window) ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	// 添加接收方
	tohbox := ui.NewHorizontalBox()
	tohbox.SetPadded(true)
	vbox.Append(tohbox, false)
	to := ui.NewEntry()
	tohbox.Append(ui.NewLabel("接收方:           "), false)
	tohbox.Append(to, true)

	// 添加金额
	amounthbox := ui.NewHorizontalBox()
	amounthbox.SetPadded(true)
	vbox.Append(amounthbox, false)
	amount := ui.NewEntry()
	amounthbox.Append(ui.NewLabel("金额:  "), false)
	amounthbox.Append(amount, true)

	// 选择币种类型
	currecyhbox := ui.NewHorizontalBox()
	currecyhbox.SetPadded(true)
	vbox.Append(currecyhbox, false)
	cbox := ui.NewCombobox()
	cbox.Append("VET")
	cbox.Append("VTHO")
	cbox.Append("ETH")
	cbox.Append("ERC20")
	currecyhbox.Append(ui.NewLabel("货币类型:"), false)
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
	currecyhbox.Append(ui.NewLabel("token地址:"), false)
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
	gashbox.Append(ui.NewLabel("Gas价格:    "), false)
	gashbox.Append(slider, true)
	gashbox.Append(spinbox, false)

	txBtn := ui.NewButton("转账")
	vbox.Append(txBtn, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	oplog := ui.NewMultilineEntry()
	oplog.SetReadOnly(true)
	vbox.Append(oplog, true)

	txBtn.OnClicked(func(*ui.Button) {
		control.Transfer(currecy, amount.Text(), to.Text(), erc20Addr.Text(), oplog)
	})
	return vbox
}
