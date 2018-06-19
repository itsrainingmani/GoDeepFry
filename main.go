package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/tsmanikandan/GoDeepFry/emoji"
	"github.com/tsmanikandan/GoDeepFry/noise"

	"github.com/disintegration/gift"
)

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

	return emoji.LoadImage(strings.Join(memeLoc, ""))
}

func main() {
	fmt.Println("Welcome to the Deep Fryer")

	randImgPtr := flag.Bool("r", false, "Random DeepFry")
	specImgPtr := flag.String("i", "", "Choose a specific image from the meme folder")
	jpegQualPtr := flag.Int("q", 100, "JPEG Image quality")
	spNoisePtr := flag.Float64("s", 0, "Amount of Salt and Pepper Noise")
	gausPtr := flag.Float64("g", 1, "Std Dev of Gaussian distribution")

	// guasNoiseImgPtr := flag.Bool("g", false, "Add Gaussian noise to a test image")

	flag.Parse()

	if *randImgPtr == true && *specImgPtr != "" {
		fmt.Println("Generating random deepfry image!")
		// rImg := randMeme()
		rImg := emoji.LoadImage(*specImgPtr)

		rand.Seed(time.Now().Unix())
		satVal := (rand.Float32() * 25) + 45
		conVal := (rand.Float32() * 20) + 50
		gamVal := (rand.Float32() * 0.8) + 1
		gaussVal := (rand.Float64()) + 0.5
		spVal := (rand.Float64() * 0.008) + 0.014
		jpegVal := rand.Intn(4) + 3

		fmt.Print(satVal, conVal, gamVal, gaussVal, spVal, jpegVal)

		g := gift.New(
			gift.Saturation(satVal),
			gift.Contrast(conVal),
			gift.Gamma(gamVal),
		)
		dst := image.NewRGBA(g.Bounds(rImg.Bounds()))

		g.Draw(dst, rImg)
		emoji.SaveImage("./deepfried/testImage.jpg", noise.SaltAndPepperNoise(*noise.GaussianNoise(*dst, gaussVal), spVal), jpegVal)
	} else if *specImgPtr != "" {
		fmt.Println("Deep Frying according to recipe")
		rImg := emoji.LoadImage(*specImgPtr)
		g := gift.New(
			gift.Saturation(60),
			gift.Contrast(50),
			gift.Gamma(1.6),
		)
		dst := image.NewRGBA(g.Bounds(rImg.Bounds()))

		g.Draw(dst, rImg)
		// saveImage("./deepfried/testImage.jpg", noise.SaltAndPepperNoise(*dst, float32(*spNoisePtr)), *jpegQualPtr)
		emoji.SaveImage("./deepfried/testImage.jpg", noise.SaltAndPepperNoise(*noise.GaussianNoise(*dst, *gausPtr), *spNoisePtr), *jpegQualPtr)
	} else {
		fmt.Println("Improper flags selected! Use the -h flag to the right usage")
	}
}
