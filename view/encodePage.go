package view

import (
	"github.com/andlabs/ui"
	"github.com/wupeaking/vechaintool/control"
)

//MakeEncodingPage 生成编码页
func MakeEncodingPage(mainwin *ui.Window) ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(false)

	grid := ui.NewGrid()
	grid.SetPadded(false)
	vbox.Append(grid, false)
	log := ui.NewMultilineEntry()
	log.SetReadOnly(false)

	// ABI 编码字符串
	strEntry := ui.NewMultilineEntry()

	abiEncodeBtn := ui.NewButton("ABI编码")
	abiEncodeBtn.OnClicked(func(*ui.Button) {
		control.ABIPack(strEntry.Text(), log)
	})

	abiDecodeBtn := ui.NewButton("ABI解码")
	abiDecodeBtn.OnClicked(func(*ui.Button) {
		control.ABIUnPack(strEntry.Text(), log)
	})

	signBtn := ui.NewButton("secp256k1签名")
	signBtn.OnClicked(func(*ui.Button) {
		control.Signature(strEntry.Text(), log)
	})

	hash256Btn := ui.NewButton("SHA256哈希")
	hash256Btn.OnClicked(func(*ui.Button) {
		control.Sha256(strEntry.Text(), log)
	})

	keccak256Btn := ui.NewButton("Keccak256哈希")
	keccak256Btn.OnClicked(func(*ui.Button) {
		control.Keccak256Hash(strEntry.Text(), log)
	})

	base64Encode := ui.NewButton("base64编码")
	base64Encode.OnClicked(func(*ui.Button) {
		control.Base64Encode(strEntry.Text(), log)
	})

	base64Decode := ui.NewButton("base64解码")
	base64Decode.OnClicked(func(*ui.Button) {
		control.Base64Decode(strEntry.Text(), log)
	})

	strWidth := 5
	grid.Append(strEntry,
		0, 0, strWidth, 7,
		true, ui.AlignFill, true, ui.AlignFill)
	grid.Append(abiEncodeBtn,
		strWidth, 0, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(abiDecodeBtn,
		strWidth, 1, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(signBtn,
		strWidth, 2, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(hash256Btn,
		strWidth, 3, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(keccak256Btn,
		strWidth, 4, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(base64Encode,
		strWidth, 5, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	grid.Append(base64Decode,
		strWidth, 6, 1, 1,
		true, ui.AlignFill, true, ui.AlignFill)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	vbox.Append(log, true)

	return vbox
}
