package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

func loadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("os.Open failed: %v", err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("image.Decode failed: %v", err)
	}
	return img
}

func saveImage(filename string, img image.Image, qual int) {
	err := imaging.Save(img, filename, imaging.JPEGQuality(qual))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

func randMeme() image.Image {
	fmt.Println("Reading from the meme folder!")
	files, err := ioutil.ReadDir("./memes")
	if err != nil {
		log.Fatalf("ioutil.ReadDir failed %v", err)
	}
	rand.Seed(time.Now().Unix())
	fmt.Println("Picking a random spicy meme!")
	randomFile := files[rand.Intn(len(files))]

	fmt.Println(randomFile.Name())
	// meme := loadImage(randomFile.Name())

	var memeLoc []string
	memeLoc = append(memeLoc, "./memes/")
	memeLoc = append(memeLoc, randomFile.Name())

	return loadImage(strings.Join(memeLoc, ""))
}

func main() {
	fmt.Println("Welcome to the Deep Fryer")

	randImgPtr := flag.Bool("r", false, "Random DeepFry")
	specImgPtr := flag.String("i", "", "Choose a specific image from the meme folder")
	jpegQualPtr := flag.Int("q", 100, "JPEG Image quality")

	// guasNoiseImgPtr := flag.Bool("g", false, "Add Gaussian noise to a test image")

	flag.Parse()

	if *randImgPtr == true && *specImgPtr == "" {
		fmt.Println("Generating random deepfry image!")
		rImg := randMeme()

		rImg = imaging.AdjustContrast(rImg, 65)
		rImg = imaging.Sharpen(rImg, 10)
		saveImage("./deepfried/testImage.jpg", rImg, *jpegQualPtr)
	} else if *specImgPtr != "" {
		fmt.Println("Deep Frying according to recipe")
		rImg := loadImage(*specImgPtr)
		rImg = imaging.AdjustContrast(rImg, 65)
		rImg = imaging.Sharpen(rImg, 10)
		saveImage("./deepfried/testImage.jpg", rImg, *jpegQualPtr)
	} else {
		fmt.Println("Improper flags selected! Use the -h flag to the right usage")
	}
}
