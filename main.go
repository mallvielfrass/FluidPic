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

func (box *FiguresCollection) AddItem(item Item) []Item {
	box.Item = append(box.Item, item)
	return box.Item
}
func DeleteItem(x, y int, box FiguresCollection) FiguresCollection {
	ItemList := []Item{}
	Localbox := FiguresCollection{ItemList}
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
func (box *FiguresCollection) FindAll(xk, yk int, m image.Image) {
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
}

func main() {
	//get original image name from args
	arg := os.Args[1]
	if arg == "" {
		fmt.Println("No argument")
		os.Exit(1)
	}
	//trim space from image name
	arg = removeSpace(arg)
	fmt.Printf("Image name:[%s]\n", arg)
	//open original image
	imgfile, err := os.Open(arg)
	if err != nil {
		fmt.Printf("imageOpen Error: %v\n", err)
		os.Exit(1)
	}
	//
	defer imgfile.Close()
	//get image width and height
	originalImgParams, err := imageParams(imgfile)
	if err != nil {
		fmt.Printf("imageParams Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nWidth : ", originalImgParams.width)
	fmt.Println("Height : ", originalImgParams.height)
	//read image data
	originalImageCanvas, _, err := image.Decode(imgfile)
	if err != nil {
		fmt.Printf("imageDecode Error: %v\n", err)
		os.Exit(1)
	}
	//create virtual image copy
	virtualImageCanvas := image.NewRGBA(image.Rect(0, 0, originalImgParams.width, originalImgParams.height))
	draw.Draw(virtualImageCanvas, virtualImageCanvas.Bounds(), originalImageCanvas, image.Point{0, 0}, draw.Src)
	ItemList := []Item{}
	box := FiguresCollection{ItemList}
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
	//	fmt.Printf("center: %d %d \n", xk, yk)
	for h := 0; h < originalImgParams.width; h++ {
		for j := 0; j < originalImgParams.height; j++ {
			// Work Zone
			//
			//
			c := check(h, j, virtualImageCanvas)
			if c != 0 {
				var xk int = h
				var yk int = j
				//c, cUp, cDown, cLeft, cRigth := getAll(xk, yk, m)
				//	fmt.Println("color=", c)

				boxDuble := FiguresCollection{ItemList}
				mcopy := image.NewRGBA(image.Rect(0, 0, originalImgParams.width, originalImgParams.height))
				s := strconv.Itoa(item)
				copyto, err := os.Create(path + "/" + s + ".png")
				if err != nil {
					log.Fatal(err)
				}
				defer copyto.Close()
				item = item + 1
				virtualImageCanvas.Set(xk, yk, color.RGBA{0, 0, 0, 255})
				box.FindAll(xk, yk, virtualImageCanvas)
				boxDuble.FindAll(xk, yk, virtualImageCanvas)
				//}
				//		fmt.Println("box:")
				fmt.Println(box)
				fmt.Println("Start cicle__________________")

				for {
					//		fmt.Println("Iteration")

					if len(box.Item) != 0 {
						for i, v := range box.Item {
							fmt.Printf("inter %d\n", i)
							c := check(v.X, v.Y, virtualImageCanvas)
							if c != 0 {
								fmt.Printf("Pixel  %d,%d is White\n", v.X, v.Y)
							} else {
								fmt.Printf("Pixel  %d,%d is Black\n", v.X, v.Y)
							}

							//			fmt.Printf("box len %d c=%d\n", len(box.Item), c)

							if c != 0 {
								fmt.Printf("Ranger c!=0 %d\n", c)
								virtualImageCanvas.Set(v.X, v.Y, color.RGBA{0, 0, 0, 255})
								fmt.Printf("Pixel  %d,%d Set as Black\n", v.X, v.Y)
								box = DeleteItem(v.X, v.Y, box)

								box.FindAll(v.X, v.Y, virtualImageCanvas)
								boxDuble.FindAll(v.X, v.Y, virtualImageCanvas)
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
