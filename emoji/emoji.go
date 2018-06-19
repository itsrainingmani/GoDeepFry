// Package emoji contains utility functions for adding emojis
// to a given image
package emoji

import (
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
)

type pos struct {
	x, y int
}

// LoadImage takes in a filename, reads the image if any from the location and returns it
func LoadImage(filename string) image.Image {
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

// SaveImage takes in a filename, an image and a quality index and saves the image to th specified location
func SaveImage(filename string, img image.Image, qual int) {
	err := imaging.Save(img, filename, imaging.JPEGQuality(qual))
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

func getAssetNames(assetDir string) []string {
	var filenames []string

	files, err := ioutil.ReadDir(assetDir)
	if err != nil {
		log.Fatalf("ioutil.ReadDir failed %v", err)
	}

	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames
}

func loadRandomAssets(assetNames []string, numToLoad int) []image.Image {
	rand.Seed(time.Now().Unix())
	randAssets := []image.Image{}
	set := make(map[int]struct{})

	for {
		randValue := rand.Intn(numToLoad)
		if _, ok := set[randValue]; ok {
			continue
		} else {
			set[randValue] = struct{}{}
		}
		if len(set) == numToLoad {
			break
		}
	}

	for i := range set {
		randAssets = append(randAssets, LoadImage(assetNames[i]))
	}

	return randAssets
}

func genRandomEmojiPositions(emjBnds image.Rectangle, srcBnds image.Rectangle) image.Pt {
	// randPos := []pos{}
	DX := srcBnds.Dx()
	DY := srcBnds.Dy()

	dx := emjBnds.Dx()
	dy := emjBnds.Dy()
	rand.Seed(time.Now().Unix())
	randX := rand.Intn(DX-dx) + dx/2
	randY := rand.Intn(DY-dy) + dy/2
	return image.Pt(randX, randY)
}

// AddEmojis is a function that will add emojis randomly to the given
// image
func AddEmojis(src image.Image) *image.RGBA {
	rand.Seed(time.Now().Unix())

	emoImg := image.NewRGBA(src.Bounds())
	gift.New().Draw(emoImg, src)

	assets := getAssetNames("../assets/")
	numAssetsToUse := rand.Intn(len(assets))

	randAssets := loadRandomAssets(assets, numAssetsToUse)

	for i := randAssets {
		randPos := genRandomEmojiPositions(randAssets[i].Bounds, src.Bounds)
		gift.New().DrawAt(emoImg, randAssets[i], randPos, gift.OverOperator)
	}

	return emoImg
}
