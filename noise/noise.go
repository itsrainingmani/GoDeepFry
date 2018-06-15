package noise

import (
	"image"
	"image/color"
	"math/rand"
	"time"
)

// SaltAndPepperNoise is a function that generates salt and pepper noise in an image
// based on a given probability
func SaltAndPepperNoise(src image.RGBA, prob float64) *image.RGBA {
	spImg := image.NewRGBA(src.Bounds())
	threshold := 1 - prob
	rand.Seed(time.Now().Unix())

	for x := 0; x < spImg.Bounds().Dx(); x++ {
		for y := 0; y < spImg.Bounds().Dy(); y++ {

			randProb := rand.Float64()
			if randProb < prob {
				spImg.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else if randProb > threshold {
				spImg.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				spImg.Set(x, y, src.At(x, y))
			}
		}
	}
	return spImg
}

// GaussianNoise is a function that adds guassian white noise to an image
func GaussianNoise(src image.RGBA, stddev float64) *image.RGBA {
	rand.Seed(time.Now().Unix())
	gImg := image.NewRGBA(src.Bounds())

	for x := 0; x < src.Bounds().Dx(); x++ {
		for y := 0; y < src.Bounds().Dy(); y++ {
			gaussProb := [3]float64{rand.NormFloat64() * stddev, rand.NormFloat64() * stddev, rand.NormFloat64() * stddev}
			r, g, b, a := src.At(x, y).RGBA()
			newR, newG, newB := uint8(float64(r)+gaussProb[0]), uint8(float64(g)+gaussProb[1]), uint8(float64(b)+gaussProb[2])
			gImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a)})
		}
	}
	return gImg
}
