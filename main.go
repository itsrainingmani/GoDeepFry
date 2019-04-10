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

	"github.com/tsmanikandan/GoDeepFry/effects"

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

	return effects.LoadImage(strings.Join(memeLoc, ""))
}

func main() {
	initString := `
    ______          ____                        ______            
   / ____/____     / __ \ ___   ___   ____     / ____/_____ __  __
  / / __ / __ \   / / / // _ \ / _ \ / __ \   / /_   / ___// / / /
 / /_/ // /_/ /  / /_/ //  __//  __// /_/ /  / __/  / /   / /_/ / 
 \____/ \____/  /_____/ \___/ \___// .___/  /_/    /_/    \__, /  
				  /_/                    /____/`
	fmt.Println(initString)
	fmt.Println("Welcome to the Deep Fryer")

	randImgPtr := flag.Bool("r", false, "Randomly generate Deep Fry")
	specImgPtr := flag.String("i", "", "Choose a specific image from the meme folder")
	jpegQualPtr := flag.Int("q", 100, "JPEG Image quality")
	spNoisePtr := flag.Float64("s", 0, "Amount of Salt and Pepper Noise")
	gausPtr := flag.Float64("g", 1, "Std Dev of Gaussian distribution")

	// guasNoiseImgPtr := flag.Bool("g", false, "Add Gaussian noise to a test image")

	flag.Parse()

	if *specImgPtr == "" {
		fmt.Println("Please specify the path of the image you want to deep fry using the -i flag")
		os.Exit(0)
	}

	if *randImgPtr == true {
		fmt.Println("Picking a random Deep Fry recipe!")
		// rImg := randMeme()
		rImg := effects.LoadImage(*specImgPtr)

		rand.Seed(time.Now().Unix())
		satVal := (rand.Float32() * 25) + 45
		conVal := (rand.Float32() * 20) + 50
		gamVal := (rand.Float32() * 0.8) + 1
		gaussVal := (rand.Float64()) + 0.5
		spVal := (rand.Float64() * 0.008) + 0.014
		jpegVal := rand.Intn(4) + 3

		fmt.Println(satVal, conVal, gamVal, gaussVal, spVal, jpegVal)

		g := gift.New(
			gift.Saturation(satVal),
			gift.Contrast(conVal),
			gift.Gamma(gamVal),
		)

		rImg = effects.AddEmojis(rImg)
		dst := image.NewRGBA(g.Bounds(rImg.Bounds()))

		g.Draw(dst, rImg)
		effects.SaveImage("./deepfried/testImage.jpg", effects.SaltAndPepperNoise(*effects.GaussianNoise(*dst, gaussVal), spVal), jpegVal)
	} else {
		fmt.Println("Deep Frying according to recipe")
		rImg := effects.LoadImage(*specImgPtr)
		g := gift.New(
			gift.Saturation(60),
			gift.Contrast(50),
			gift.Gamma(1.6),
		)
		dst := image.NewRGBA(g.Bounds(rImg.Bounds()))

		g.Draw(dst, rImg)
		effects.SaveImage("./deepfried/testImage.jpg", effects.SaltAndPepperNoise(*effects.GaussianNoise(*dst, *gausPtr), *spNoisePtr), *jpegQualPtr)
	}
}
