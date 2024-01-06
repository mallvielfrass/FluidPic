package main

import (
	"fmt"
	"image"
	"os"
	"unicode"
)

func removeSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

func imageParams(img *os.File) (ImageParams, error) {
	imgCfg, _, err := image.DecodeConfig(img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	img.Seek(0, 0) //return reader to begin of file
	width := imgCfg.Width
	height := imgCfg.Height
	return ImageParams{width, height}, nil
}
