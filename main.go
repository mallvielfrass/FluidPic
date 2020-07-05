package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
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
			//	fmt.Printf("%d %d deleted\n", x, y)
		} else {
			Localbox.AddItem(box.Item[i])
		}
	}
	return Localbox
}
func FindAll(box MainBox, xk, yk int, m image.Image) MainBox {
	_, cUp, cDown, cLeft, cRigth := getAll(xk, yk, m)
	if cUp != 0 {
		if 0 <= (yk - 1) {
			//	fmt.Println("cUp != 0 ")
			box.AddItem(Item{
				X: xk,
				Y: yk - 1,
			})
		}
	}
	if cDown != 0 {
		//	fmt.Println("cDown != 0 ")
		box.AddItem(Item{
			X: xk,
			Y: yk + 1,
		})
	}
	if cLeft != 0 {
		if 0 <= (xk - 1) {
			//	fmt.Println("cLeft != 0 ")
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
	arg := os.Args[1]
	//printl(arg)
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
	img, _, err := image.Decode(imgfile)
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), img, image.Point{0, 0}, draw.Src)
	ItemList := []Item{}
	box := MainBox{ItemList}
	var xk int = (width / 2)
	var yk int = (height / 2)
	item := 0
	pathx := strings.Split(arg, ".")
	path := pathx[0]
	errDir := os.MkdirAll(path, 0755)
	if errDir != nil {
		panic(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
	}
	fmt.Printf("center: %d %d \n", xk, yk)
	for h := 0; h < width; h++ {
		for j := 0; j < height; j++ {
			// Work Zone
			//
			//
			c := check(h, j, m)
			if c != 0 {
				var xk int = h
				var yk int = j
				//c, cUp, cDown, cLeft, cRigth := getAll(xk, yk, m)
				//	fmt.Println("color=", c)

				boxDuble := MainBox{ItemList}
				mcopy := image.NewRGBA(image.Rect(0, 0, width, height))
				s := strconv.Itoa(item)
				copyto, err := os.Create(path + "/" + s + ".png")
				if err != nil {
					log.Fatal(err)
				}
				defer copyto.Close()
				item = item + 1
				m.Set(xk, yk, color.RGBA{0, 0, 0, 255})
				//fmt.Println("4sleep:")
				//		fmt.Println(box)
				//time.Sleep(4000 * time.Millisecond)
				if c != 0 {
					//		fmt.Printf("c!=0 %d\n", c)
					m.Set(xk, yk, color.RGBA{0, 0, 0, 255})
					box = FindAll(box, xk, yk, m)
					boxDuble = FindAll(boxDuble, xk, yk, m)
				}
				//		fmt.Println("box:")
				fmt.Println(box)
				fmt.Println("Start cicle__________________")

				for {
					//		fmt.Println("Iteration")

					if len(box.Item) != 0 {
						for i, v := range box.Item {
							fmt.Printf("inter %d\n", i)
							c := check(v.X, v.Y, m)
							if c != 0 {
								fmt.Printf("Pixel  %d,%d is White\n", v.X, v.Y)
							} else {
								fmt.Printf("Pixel  %d,%d is Black\n", v.X, v.Y)
							}

							//			fmt.Printf("box len %d c=%d\n", len(box.Item), c)

							if c != 0 {
								fmt.Printf("Ranger c!=0 %d\n", c)
								m.Set(v.X, v.Y, color.RGBA{0, 0, 0, 255})
								fmt.Printf("Pixel  %d,%d Set as Black\n", v.X, v.Y)
								box = DeleteItem(v.X, v.Y, box)

								box = FindAll(box, v.X, v.Y, m)
								boxDuble = FindAll(boxDuble, v.X, v.Y, m)
							} else {
								fmt.Printf("Pixel  %d,%d already is Black\n", v.X, v.Y)
								box = DeleteItem(v.X, v.Y, box)
							}
							//			fmt.Println("box in iter:")
							//			fmt.Println(box)
							//time.Sleep(2000 * time.Millisecond)
						}
					} else {
						break
					}
				}
				//	fmt.Print("box: ")
				//		fmt.Println(box)
				//		fmt.Println(boxDuble)
				//Work Zone Closed
				for {
					fmt.Println("Iteration output create")

					if len(boxDuble.Item) != 0 {
						for _, z := range boxDuble.Item {

							mcopy.Set(z.X, z.Y, color.RGBA{0, 0, 0, 255})
							fmt.Printf("BoxSet %d,%d Set as Black\n", z.X, z.Y)
							boxDuble = DeleteItem(z.X, z.Y, boxDuble)

						}
					} else {
						break
					}
				}
				png.Encode(copyto, mcopy)
			}
		}
	}
	//png.Encode(outFile, m)

}
