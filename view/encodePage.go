package view

import (
	"github.com/andlabs/ui"
	"github.com/wupeaking/vechaintool/control"
)

//MakeEncodingPage 生成编码页
func MakeEncodingPage(mainwin *ui.Window) ui.Control {
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

	signBtn := ui.NewButton("secp256k1 signature")
	signBtn.OnClicked(func(*ui.Button) {
		control.Signature(strEntry.Text(), log)
	})

	hash256Btn := ui.NewButton("SHA256")
	hash256Btn.OnClicked(func(*ui.Button) {
		control.Sha256(strEntry.Text(), log)
	})

	keccak256Btn := ui.NewButton("Keccak256")
	keccak256Btn.OnClicked(func(*ui.Button) {
		control.Keccak256Hash(strEntry.Text(), log)
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

	vbox.Append(log, true)

	return vbox
}
