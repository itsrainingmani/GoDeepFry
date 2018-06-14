package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
)

// loadImage takes in a filename, reads the image if any from the location and returns it
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

// saveImage takes in a filename, an image and a quality index and saves the image to th specified location
func saveImage(filename string, img image.Image, qual int) {
	err := imaging.Save(img, filename, imaging.JPEGQuality(qual))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

// randMeme returns a random meme image from the memes folder
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

	// We need to construct a string array so that we can join the contents to make
	// a filename that includes the folder location
	var memeLoc []string
	memeLoc = append(memeLoc, "./memes/")
	memeLoc = append(memeLoc, randomFile.Name())

	return loadImage(strings.Join(memeLoc, ""))
}

// func makeFilters(satInc float32, contInc float32, gamInc float32) *gift.GIFT {
// 	return gift.New(
// 		gift.Saturation(satInc),
// 		gift.Contrast(contInc),
// 		gift.Gamma(gamInc),
// 	)
// }

func saltAndPepperNoise(src image.RGBA, prob float32) *image.RGBA {
	spImg := image.NewRGBA(src.Bounds())
	threshold := 1 - prob
	rand.Seed(time.Now().Unix())

	for x := 0; x < spImg.Bounds().Dx(); x++ {
		for y := 0; y < spImg.Bounds().Dy(); y++ {

			randProb := rand.Float32()
			if randProb < prob {
				spImg.Set(x, y, color.RGBA{0, 0, 0, 0})
			} else if randProb > threshold {
				spImg.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else {
				spImg.Set(x, y, src.At(x, y))
			}
		}
	}
	return spImg
}

func main() {
	fmt.Println("Welcome to the Deep Fryer")

	randImgPtr := flag.Bool("r", false, "Random DeepFry")
	specImgPtr := flag.String("i", "", "Choose a specific image from the meme folder")
	jpegQualPtr := flag.Int("q", 100, "JPEG Image quality")
	spNoisePtr := flag.Float64("s", 0, "Amount of Salt and Pepper Noise")

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
		g := gift.New(
			gift.Saturation(50),
			gift.Contrast(40),
			gift.Gamma(1.5),
		)
		dst := image.NewRGBA(g.Bounds(rImg.Bounds()))

		g.Draw(dst, rImg)
		saveImage("./deepfried/testImage.jpg", saltAndPepperNoise(*dst, float32(*spNoisePtr)), *jpegQualPtr)
	} else {
		fmt.Println("Improper flags selected! Use the -h flag to the right usage")
	}
}
