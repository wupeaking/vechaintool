package widget

import "github.com/andlabs/ui"

func ConfirmDialog(msg string, confirmCb func(), cancelCb func()) {
	w := ui.NewWindow("Confirm", 300, 100, true)
	w.SetMargined(true)
	w.OnClosing(func(*ui.Window) bool {
		w.Destroy()
		cancelCb()
		return true
	})

	confirm := ui.NewButton("确定")
	confirm.OnClicked(func(b *ui.Button) {
		w.Destroy()
		confirmCb() // or go callback() depending on what it does
	})

	cancel := ui.NewButton("取消")
	cancel.OnClicked(func(b *ui.Button) {
		w.Destroy()
		cancelCb()
	})

	wrapper := ui.NewVerticalBox()
	wrapper.SetPadded(true)
	wrapper.Append(ui.NewLabel(msg), false)
	wrapper.Append(confirm, false)
	wrapper.Append(cancel, false)
	w.SetChild(wrapper)
	w.Show()

}
