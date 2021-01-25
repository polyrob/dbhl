package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	jpeg "image/jpeg"
	png "image/png"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

func main() {
	profileUrl := os.Args[1]
	fmt.Printf("Start with source image, %s\n", profileUrl)

	// load template
	templateFile, err := os.Open("./template.png")
	if err != nil {
		log.Fatalln("Could not open template", err)
	}
	defer templateFile.Close()

	templateImage, err := png.Decode(templateFile)
	if err != nil {
		log.Fatalln("Could not decode image file", err)
	}
	width := templateImage.Bounds().Max.X
	height := templateImage.Bounds().Max.Y
	fmt.Printf("Template Image width: %d, height: %d\n", width, height)

	// download remote image from url
	response, err := http.Get(profileUrl)
	if err != nil {
		log.Fatalln("Could not download remote image", err)
	}

	defer response.Body.Close()
	profileImage, err := png.Decode(response.Body)
	if err != nil {
		log.Fatalln("Could not decode remote image file", err)
	}
	profileImage = imaging.Resize(profileImage, 600, 0, imaging.Lanczos)
	profileImage = imaging.Rotate(profileImage, 42.0, color.Transparent)

	// how to resize: dstImage2 := imaging.Resize(src2, 256, 256, imaging.Lanczos)
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	dst := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// draw profile
	draw.Draw(dst, profileImage.Bounds(), profileImage, image.Point{20, 260}, draw.Src)

	// draw template over it
	draw.Draw(dst, templateImage.Bounds(), templateImage, image.ZP, draw.Over)

	out, err := os.Create("./output.jpg")
	if err != nil {
		log.Fatalln("Could create file", err)
	}

	var opt jpeg.Options
	opt.Quality = 80
	jpeg.Encode(out, dst, &opt)

	fmt.Println("End")
}
