package main

import (
	"image"
)

func checkColor(xk, yk int, img image.Image) int {
	rf, gf, bf, _ := img.At(xk, yk).RGBA()
	R := int(float32(rf) * (255.0 / 65535.0))
	G := int(float32(gf) * (255.0 / 65535.0))
	B := int(float32(bf) * (255.0 / 65535.0))
	Color := (R + G + B) / 3
	// if Color == 0 {
	// 	fmt.Printf("color:= Black\n")
	// }
	// if Color == 1 {
	// 	fmt.Printf("color:= White\n")
	// }
	if Color == 255 {
		return 1
	}
	return 0

}
func getNeighbours(xk, yk int, img image.Image) (c, cUp, cDown, cLeft, cRigth int) {
	c = checkColor(xk, yk, img)
	cUp = checkColor(xk, yk-1, img)
	cDown = checkColor(xk, yk+1, img)
	cLeft = checkColor(xk-1, yk, img)
	cRigth = checkColor(xk+1, yk, img)
	return c, cUp, cDown, cLeft, cRigth
}
