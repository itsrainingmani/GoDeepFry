package noise

import (
	"image"
	"image/color"
	"math/rand"
	"time"
)

// SaltAndPepperNoise is a function that generates salt and pepper noise in an image
// based on a given probability
func SaltAndPepperNoise(src image.RGBA, prob float32) *image.RGBA {
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
