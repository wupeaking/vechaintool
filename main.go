package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/wupeaking/vechaintool/view"
	// "github.com/ethereum/go-ethereum/crypto"
)

var mainwin *ui.Window

func setupUI() {
	mainwin = ui.NewWindow("区块链智能合约调试助手 V1.0", 640, 480, true)
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

	tab.Append("设置", view.MakeSettingPage(mainwin))
	tab.SetMargined(0, true)

	tab.Append("交易", view.MakeTransferPage(mainwin))
	tab.SetMargined(1, true)

	tab.Append("合约调试", view.MakeContractPage(mainwin))
	tab.SetMargined(2, true)

	tab.Append("编码", view.MakeEncodingPage(mainwin))
	tab.SetMargined(3, true)
	view.RegistRefreshPage("settingPage", func() {
		tab.Delete(2)
		tab.InsertAt("合约调试", 2, view.MakeContractPage(mainwin))
	})

	go view.StartRefresh()

	//go func() {
	//	for {
	//		select {
	//		case <-refresh:
	//			ui.QueueMain(func() {
	//				tab.Delete(2)
	//				tab.InsertAt("contract", 2, view.MakeContractPage(mainwin))
	//			})
	//		}
	//	}
	//}()
	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
