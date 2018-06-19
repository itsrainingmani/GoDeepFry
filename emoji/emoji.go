// Package emoji contains utility functions for adding emojis
// to a given image
package emoji

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
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
		filenames = append(filenames, strings.Join([]string{assetDir, file.Name()}, ""))
	}

	return filenames
}

func loadRandomAssets(assetNames []string, numToLoad int) []image.Image {

	randAssets := []image.Image{}

	// Since Golang does not come with pre-built support for a set type this is
	// one possible workaround. This creates a map of int to an empty struct.
	// when a generate a random value, we check if that key has a struct, if it does
	// then the element exists. if not we add a new struct. This is more space efficient
	// than making a map of int to bool.
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

func loadMultipleRandomAssets(assetNames []string) []image.Image {
	randAssets := []image.Image{}

	assetMap := make(map[string]int)

	for _, asset := range assetNames {
		randValue := rand.Intn(2) + 1
		assetMap[asset] = randValue
	}

	fmt.Println(assetMap)

	for k, v := range assetMap {
		if v != 0 {
			for i := 0; i < v; i++ {
				randAssets = append(randAssets, LoadImage(k))
			}
		}
	}

	return randAssets
}

func genRandomEmojiPositions(emjBnds image.Rectangle, srcBnds image.Rectangle) image.Point {

	DX := srcBnds.Dx()
	DY := srcBnds.Dy()

	dx := emjBnds.Dx()
	dy := emjBnds.Dy()

	// This is to ensure that the emoji's position does not cross the edges of the source image
	randX := rand.Intn(DX-dx) + dx/2
	randY := rand.Intn(DY-dy) + dy/2
	return image.Pt(randX, randY)
}

// AddEmojis is a function that will add emojis randomly to the given
// image
func AddEmojis(src image.Image) *image.RGBA {
	rand.Seed(time.Now().Unix())

	// Makes a new empty image and then copies over the source
	emoImg := image.NewRGBA(src.Bounds())
	gift.New().Draw(emoImg, src)

	fmt.Println("Getting Asset names")
	assets := getAssetNames("./assets/")
	fmt.Println(assets)
	numAssetsToUse := rand.Intn(len(assets)-1) + 1

	fmt.Println("Number of assets to use - ", numAssetsToUse)

	fmt.Println("Loading Random assets!")
	// loadMultipleRandomAssets(assets)
	// randAssets := loadRandomAssets(assets, numAssetsToUse)
	randAssets := loadMultipleRandomAssets(assets)

	// Loop over the random assets and add them to the new image at random positions
	fmt.Println("Adding random emojis to the image!")
	for _, rAss := range randAssets {
		randPos := genRandomEmojiPositions(rAss.Bounds(), src.Bounds())
		gift.New().DrawAt(emoImg, rAss, randPos, gift.OverOperator)
	}

	return emoImg
}
