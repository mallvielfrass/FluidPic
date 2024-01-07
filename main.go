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

// main.go
func init() {
	// damn important or else At(), Bounds() functions will
	// caused memory pointer error!!
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

func (box *FiguresCollection) AddPoint(item Point) []Point {
	box.Points = append(box.Points, item)
	return box.Points
}

func (box *FiguresCollection) AddNeighbours(xk, yk int, m image.Image) {
	_, cUp, cDown, cLeft, cRigth := getNeighbours(xk, yk, m)
	if cUp != 0 {
		if 0 <= (yk - 1) {
			//	fmt.Println("cUp != 0 ")
			box.AddPoint(Point{
				X: xk,
				Y: yk - 1,
			})
		}
	}
	if cDown != 0 {
		//	fmt.Println("cDown != 0 ")
		box.AddPoint(Point{
			X: xk,
			Y: yk + 1,
		})
	}
	if cLeft != 0 {
		if 0 <= (xk - 1) {
			//	fmt.Println("cLeft != 0 ")
			box.AddPoint(Point{
				X: xk - 1,
				Y: yk,
			})
		}
	}
	if cRigth != 0 {
		box.AddPoint(Point{
			X: xk + 1,
			Y: yk,
		})
	}
}
func fromAgrsOrDefauilt(args []string, def string) string {
	if 2 <= len(args) {
		def = args[1]
	}
	return removeSpace(def)
}
func main() {
	//get original image name from args

	arg := fromAgrsOrDefauilt(os.Args, "origin.png")
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
	PointList := []Point{}
	box := FiguresCollection{PointList}
	item := 0
	pathx := strings.Split(arg, ".")
	path := pathx[0]
	errDir := os.MkdirAll(path, 0755)
	if errDir != nil {
		panic(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal(err)

	}
	fmt.Println("start")
	fmt.Printf("width: %d height: %d \n", originalImgParams.width, originalImgParams.height)
	//	fmt.Printf("center: %d %d \n", xk, yk)
	for h := 0; h < originalImgParams.width; h++ {
		for j := 0; j < originalImgParams.height; j++ {

			c := checkColor(h, j, virtualImageCanvas)
			if c == 0 {
				continue
			}

			var xk int = h
			var yk int = j
			//c, cUp, cDown, cLeft, cRigth := getAll(xk, yk, m)
			//	fmt.Println("color=", c)

			boxDuble := FiguresCollection{[]Point{}}
			mcopy := image.NewRGBA(image.Rect(0, 0, originalImgParams.width, originalImgParams.height))
			s := strconv.Itoa(item)
			copyto, err := os.Create(path + "/" + s + ".png")
			if err != nil {
				log.Fatal(err)
			}
			defer copyto.Close()
			item = item + 1
			virtualImageCanvas.Set(xk, yk, color.RGBA{0, 0, 0, 255})
			box.AddNeighbours(xk, yk, virtualImageCanvas)
			boxDuble.AddNeighbours(xk, yk, virtualImageCanvas)
			//}
			//		fmt.Println("box:")
			//fmt.Println(box)
			//	fmt.Println("Start cicle__________________")

			for {
				//		fmt.Println("Iteration")

				if len(box.Points) == 0 {
					break
				}
				for _, v := range box.Points {
					//	fmt.Printf("inter %d\n", i)
					box = DeletePoint(v.X, v.Y, box)
					c := checkColor(v.X, v.Y, virtualImageCanvas)
					if c == 0 {
						continue
					}
					virtualImageCanvas.Set(v.X, v.Y, color.RGBA{0, 0, 0, 255})
					box.AddNeighbours(v.X, v.Y, virtualImageCanvas)
					boxDuble.AddNeighbours(v.X, v.Y, virtualImageCanvas)
				}

			}
			//Work Zone Closed
			for {
				fmt.Println("Iteration output create")
				if len(boxDuble.Points) == 0 {
					break
				}
				for _, z := range boxDuble.Points {
					mcopy.Set(z.X, z.Y, color.RGBA{0, 0, 0, 255})
					//	fmt.Printf("BoxSet %d,%d Set as Black\n", z.X, z.Y)
					boxDuble = DeletePoint(z.X, z.Y, boxDuble)

				}
			}
			png.Encode(copyto, mcopy)
		}

	}
	fmt.Println("stop")

}
