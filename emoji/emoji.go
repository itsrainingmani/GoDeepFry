// Package emoji contains utility functions for adding emojis
// to a given image
package emoji

import (
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

// type resizeby struct {
// 	width    int
// 	height   int
// 	resample gift.Resampling
// }

// type rotateby struct {
// 	angle      float32
// 	background color.Color
// 	interpol   gift.Interpolation
// }
// type emoji struct {
// 	asset   image.Image
// 	randpos image.Point
// 	resize  resizeby
// 	rotate  rotateby
// }

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
				loadedImg := LoadImage(k)
				moddedImg := rotateAndResizeEmoji(loadedImg)
				randAssets = append(randAssets, moddedImg)
			}
		}
	}

	return randAssets
}

func genRandomEmojiPos(emjBnds image.Rectangle, srcBnds image.Rectangle) image.Point {

	// Random emoji positions
	DX := srcBnds.Dx()
	DY := srcBnds.Dy()

	dx := emjBnds.Dx()
	dy := emjBnds.Dy()

	// This is to ensure that the emoji's position does not cross the edges of the source image
	randX := rand.Intn(DX-dx) + dx
	randY := rand.Intn(DY-dy) + dy
	return image.Pt(randX, randY)
}

func rotateAndResizeEmoji(emj image.Image) image.Image {

	origWidth := emj.Bounds().Dx()
	origHeight := emj.Bounds().Dy()

	randXPerc := (rand.Float32() * 0.1) + 0.4
	randYPerc := (rand.Float32() * 0.1) + 0.4

	newResizeHeight := int(float32(origHeight) * randXPerc)
	newResizeWidth := int(float32(origWidth) * randYPerc)

	var randAngle float32
	randOrientation := rand.Float32()

	if randOrientation < 0.5 {
		randAngle = (rand.Float32() * 10.0) + 23.0
	} else {
		randAngle = (rand.Float32() * 10.0) + 330.0
	}

	g := gift.New(
		gift.Resize(newResizeHeight, newResizeWidth, gift.LanczosResampling),
		gift.Rotate(randAngle, color.Transparent, gift.CubicInterpolation),
	)
	dst := image.NewRGBA(g.Bounds(emj.Bounds()))
	g.Draw(dst, emj)

	return dst
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

	fmt.Println("Loading Random assets!")
	// loadMultipleRandomAssets(assets)
	// randAssets := loadRandomAssets(assets, numAssetsToUse)
	randAssets := loadMultipleRandomAssets(assets)

	// Loop over the random assets and add them to the new image at random positions
	fmt.Println("Adding random emojis to the image!")
	for _, rAss := range randAssets {
		randPos := genRandomEmojiPos(rAss.Bounds(), src.Bounds())
		gift.New().DrawAt(emoImg, rAss, randPos, gift.OverOperator)
	}

	return emoImg
}
