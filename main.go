package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/nfnt/resize"
)

func main() {
	var inputFile, outputDir, maxScale = parseFlags()
	r, _ := regexp.Compile(`\/(\w*)\.`)
	var prefix, _, _ = strings.Cut(r.FindString(inputFile), ".")
	var img = load(inputFile)
	var height, width int = img.Bounds().Dy(), img.Bounds().Dx()
	if math.Min(float64(height), float64(width))/float64(maxScale) >= 1 {
		for i := 2; i < int(maxScale+2); i++ {
			resized := resize.Resize(uint(width/i), 0, img, resize.Lanczos3)
			save(fmt.Sprintf("%s%s_%d.jpg", outputDir, prefix, i-1), resized) // prefix already has /
		}
	} else {
		log.Println("Image is too small")
	}

}

func parseFlags() (inputFile, outputDir string, maxScale int) {
	flag.StringVar(&inputFile, "i", "", "input file path")
	flag.StringVar(&outputDir, "o", "./", "output directory")
	flag.IntVar(&maxScale, "s", 7, "number of resized images")
	flag.Parse()
	if inputFile == "" {
		log.Fatalf("input file must be specified")
	}
	return
}

func load(filePath string) *image.RGBA {
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Cannot read file:", err)
	}
	defer imgFile.Close()

	var img image.Image
	var _, ext, _ = strings.Cut(filePath, ".")
	if ext == "png" {
		img, err = png.Decode(imgFile)
	} else {
		img, err = jpeg.Decode(imgFile)
	}
	if err != nil {
		log.Println("Cannot decode file:", err)
	}
	return img.(*image.RGBA)
}

func save(filePath string, img image.Image) {
	imgFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	defer imgFile.Close()
	jpeg.Encode(imgFile, img, nil)
}
