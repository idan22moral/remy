package lib

import "github.com/skip2/go-qrcode"

func PrintQR(qrcode *qrcode.QRCode) {
	const WHITE_SPACE string = "\033[47m  \033[0m"
	const BLACK_SPACE string = "\033[40m  \033[0m"

	for _, bitRow := range qrcode.Bitmap() {
		for _, bit := range bitRow {
			if bit {
				print(WHITE_SPACE)
			} else {
				print(BLACK_SPACE)
			}
		}
		print("\n")
	}
}
