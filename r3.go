package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func init() {
	// damn important or else At(), Bounds() functions will
	// caused memory pointer error!!
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}
func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Point struct {
	X, Y int
}

func check(xk, yk int, img image.Image) int {
	rf, gf, bf, _ := img.At(xk, yk).RGBA()
	R := int(float32(rf) * (255.0 / 65535.0))
	G := int(float32(gf) * (255.0 / 65535.0))
	B := int(float32(bf) * (255.0 / 65535.0))
	Color := (R + G + B) / 3
	if Color == 0 {
		fmt.Printf("color:= Black\n")
	}
	if Color == 1 {
		fmt.Printf("color:= White\n")
	}
	if Color == 255 {
		return 1
	}
	return 0

}
func getAll(xk, yk int, img image.Image) (c, cUp, cDown, cLeft, cRigth int) {
	c = check(xk, yk, img)
	cUp = check(xk, yk-1, img)
	cDown = check(xk, yk+1, img)
	cLeft = check(xk-1, yk, img)
	cRigth = check(xk+1, yk, img)
	return c, cUp, cDown, cLeft, cRigth
}

type MainBox struct {
	Item []Item
}
type Item struct {
	X int
	Y int
}

func (box *MainBox) AddItem(item Item) []Item {
	box.Item = append(box.Item, item)
	return box.Item
}
func DeleteItem(x, y int, box MainBox) MainBox {
	ItemList := []Item{}
	Localbox := MainBox{ItemList}
	fmt.Println(len(box.Item))
	ln := len(box.Item)
	for i := 0; i < ln; i++ {
		if box.Item[i].X == x && box.Item[i].Y == y {
			fmt.Printf("%d %d deleted\n", x, y)
		} else {
			Localbox.AddItem(box.Item[i])
		}
	}
	return Localbox
}
func FindAll(box MainBox, xk, yk int, m image.Image) MainBox {
	_, cUp, cDown, cLeft, cRigth := getAll(xk, yk, m)
	if cUp != 0 {
		if 0 < (yk - 1) {
			fmt.Println("cUp != 0 ")
			box.AddItem(Item{
				X: xk,
				Y: yk - 1,
			})
		}
	}
	if cDown != 0 {
		fmt.Println("cDown != 0 ")
		box.AddItem(Item{
			X: xk,
			Y: yk + 1,
		})
	}
	if cLeft != 0 {
		if 0 < (xk - 1) {
			fmt.Println("cLeft != 0 ")
			box.AddItem(Item{
				X: xk - 1,
				Y: yk,
			})
		}
	}
	if cRigth != 0 {
		box.AddItem(Item{
			X: xk + 1,
			Y: yk,
		})
	}
	return box
}

func main() {
	//arg := os.Args[1]
	arg := "si.png"
	printl(arg)
	imgfile, err := os.Open(arg)
	errCheck(err)

	defer imgfile.Close()
	imgCfg, _, err := image.DecodeConfig(imgfile)
	errCheck(err)

	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("\nWidth : ", width)
	fmt.Println("Height : ", height)
	imgfile.Seek(0, 0)
	var xk int = (width / 2)
	var yk int = (height / 2)
	fmt.Printf("center: %d %d \n", xk, yk)
	img, _, err := image.Decode(imgfile)

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	outFile, err := os.Create("1" + arg)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	draw.Draw(m, m.Bounds(), img, image.Point{0, 0}, draw.Src)
	// Work Zone
	//
	//
	c := check(xk, yk, m)
	//c, cUp, cDown, cLeft, cRigth := getAll(xk, yk, m)
	fmt.Println("color=", c)
	ItemList := []Item{}
	box := MainBox{ItemList}
	m.Set(xk, yk, color.RGBA{0, 0, 0, 255})
	fmt.Println("4sleep:")
	fmt.Println(box)
	//time.Sleep(4000 * time.Millisecond)
	if c != 0 {
		fmt.Printf("c!=0 %d\n", c)
		m.Set(xk, yk, color.RGBA{0, 0, 0, 255})
		box = FindAll(box, xk, yk, m)
	}
	fmt.Println("box:")
	fmt.Println(box)
	for {
		fmt.Println("Iteration")

		if len(box.Item) != 0 {
			for _, v := range box.Item {
				fmt.Printf("Debug\n")
				fmt.Printf("box len %d\n", len(box.Item))
				c := check(v.X, v.X, m)
				if c != 0 {
					fmt.Printf("Ranger c!=0 %d\n", c)
					m.Set(v.X, v.Y, color.RGBA{0, 0, 0, 255})
					box = DeleteItem(v.X, v.Y, box)

					box = FindAll(box, v.X, v.Y, m)
				} else {
					box = DeleteItem(v.X, v.Y, box)
				}
				fmt.Println("box in iter:")
				fmt.Println(box)
				//time.Sleep(2000 * time.Millisecond)
			}
		} else {
			break
		}
	}
	fmt.Println("box:")
	fmt.Println(box)
	//Work Zone Closed
	png.Encode(outFile, m)
}
