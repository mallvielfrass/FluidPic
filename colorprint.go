package main

import (
	"fmt"
	"image"
	"os"

	cr "github.com/fatih/color"
)

func intlen(x int) int {
	z := 0
	if x < 10 {
		z = 1
	}
	if x >= 10 {
		z = 2
	}
	if x > 99 {
		z = 3
	}
	return z
}
func sw(x int) string {
	sx := ""
	switch intlen(x) {
	case 1:
		sx = "  "
	case 2:
		sx = " "
	case 3:
		sx = ""
	}
	return sx
}
func printl(arg string) {
	imgfile, err := os.Open(arg)
	errCheck(err)

	defer imgfile.Close()
	imgCfg, _, err := image.DecodeConfig(imgfile)
	errCheck(err)

	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)
	imgfile.Seek(0, 0)
	img, _, err := image.Decode(imgfile)
	n := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rf, gf, bf, af := img.At(x, y).RGBA()
			R := int(float32(rf) * (255.0 / 65535.0))
			G := int(float32(gf) * (255.0 / 65535.0))
			B := int(float32(bf) * (255.0 / 65535.0))
			A := int(float32(af) * (255.0 / 65535.0))

			switch n {
			case 0:
				cr.Set(cr.FgCyan)
			case 1:
				cr.Set(cr.FgGreen)

			case 2:
				cr.Set(cr.FgYellow)

			}

			strng := fmt.Sprintf("[X : %s%d, Y : %s%d] R : %s%d, G : %s%d, B : %s%d, A : %s%d |		", sw(x), x, sw(y), y, sw(R), R, sw(G), G, sw(B), B, sw(A), A)
			fmt.Print(strng)
			cr.Unset()
			if n == 2 {
				fmt.Printf("\n")
				n = 0
			} else {
				n++
			}

		}
	}
}
