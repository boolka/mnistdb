// mnistcli is helper tool to work with "the mnist database"
// For additional info look at https://yann.lecun.com/exdb/mnist/
//
// cli arguments:
//
//	--train to use train (default) dataset (train-images.idx3-ubyte + train-labels.idx1-ubyte)
//	--test to use test dataset (t10k-images.idx3-ubyte + t10k-labels.idx1-ubyte)
//	--out (-o) chose directory to output images (make sure the directory exists)
//	--extract (-e) extract train or test database
//	--index (-i) chose image by index to extract
//	--help (-h) to print this message
//
// Suppose we want to extract all test images to "test_out" directory:
//
//	./mnistcli -e --test -o test_out
//
// Or we want to extract image by index 1000 from train db to "train_out" directory:
//
//	./mnistcli -e --train -o train_out -i 1000
//
// And so on
package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/boolka/mnistidx/pkg/mnistidx"
)

//go:embed db/train-images.idx3-ubyte
var trainImages []byte

//go:embed db/train-labels.idx1-ubyte
var trainLabels []byte

//go:embed db/t10k-images.idx3-ubyte
var testImages []byte

//go:embed db/t10k-labels.idx1-ubyte
var testLabels []byte

//go:embed help.txt
var helpTxt string

func saveDigit(filePath string, buf []byte, width, height int) {
	f, err := os.Create(filePath)

	if err != nil {
		log.Fatalln(err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			img.Set(x, y, color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: buf[y*width+x],
			})
		}
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatalln(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	args := os.Args[1:]
	var isTrain = true
	var isExtract = false
	var imageIndex int = -1
	var outputDir string = ""

	var isArgValue = false

	for i, arg := range args {
		switch arg {
		case "-h", "--help":
			fmt.Println(helpTxt)
			return
		case "-i", "--index":
			i, err := strconv.Atoi(args[i+1])

			if err != nil {
				log.Fatalln("number not recognized:", err)
			}
			isArgValue = true

			imageIndex = i
		case "--extract", "-e":
			isExtract = true
		case "--train":
			isTrain = true
		case "--test":
			isTrain = false
		case "-o", "--out":
			outputDir = args[i+1]

			isArgValue = true
		default:
			if isArgValue {
				isArgValue = false
				continue
			}

			log.Fatalln("unrecognized cli argument:", arg)
		}
	}

	var imagesBuf *bytes.Reader
	var labelsBuf *bytes.Reader

	if isTrain {
		imagesBuf = bytes.NewReader(trainImages)
		labelsBuf = bytes.NewReader(trainLabels)
	} else {
		imagesBuf = bytes.NewReader(testImages)
		labelsBuf = bytes.NewReader(testLabels)
	}

	mnistdb, err := mnistidx.NewIDX(imagesBuf, labelsBuf)

	if err != nil {
		log.Fatalln("mnistidx creation err:", err)
	}

	buf := make([]byte, mnistdb.ImageBufSize())

	if imageIndex != -1 {
		for i := 0; i <= imageIndex; i++ {
			label, err := mnistdb.Read(buf)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln(err)
			}

			if i != imageIndex {
				continue
			}

			width := int(mnistdb.ImagesHeader.ImgCols)
			height := int(mnistdb.ImagesHeader.ImgRows)

			saveDigit(filepath.Join(outputDir, fmt.Sprintf("label_%d_index_%d.png", label, i)), buf, width, height)
		}
	} else if isExtract {
		i := 0

		for {
			label, err := mnistdb.Read(buf)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln(err)
			}

			if imageIndex == -1 || i == imageIndex {
				width := int(mnistdb.ImagesHeader.ImgCols)
				height := int(mnistdb.ImagesHeader.ImgRows)

				saveDigit(filepath.Join(outputDir, fmt.Sprintf("label_%d_index_%d.png", label, i)), buf, width, height)
			}

			i++
		}

		db := "train"

		if isTrain {
			db = "test"
		}

		fmt.Println("completely extract " + db + "mnist data base")
	} else {
		fmt.Println(helpTxt)
	}
}
