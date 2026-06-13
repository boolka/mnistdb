package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"

	db "github.com/boolka/mnistdb/pkg/mnistdb"
	"github.com/boolka/mnistidx/pkg/mnistidx"
)

//go:embed help.txt
var helpTxt string

func main() {
	var isHelp, isTrain, isTest, isPng bool
	var imageIndex int
	var outputDir string

	flag.BoolVar(&isHelp, "help", false, "")
	flag.BoolVar(&isHelp, "h", false, "")
	flag.BoolVar(&isTrain, "train", false, "")
	flag.BoolVar(&isTest, "test", false, "")
	flag.BoolVar(&isPng, "png", false, "")
	flag.IntVar(&imageIndex, "index", -1, "")
	flag.IntVar(&imageIndex, "i", -1, "")
	flag.StringVar(&outputDir, "dir", "", "")
	flag.StringVar(&outputDir, "d", "", "")

	flag.Parse()

	if isHelp {
		log.Println(helpTxt)
		return
	}

	if isPng {
		var imagesBuf *bytes.Reader
		var labelsBuf *bytes.Reader

		// if train set or both omitted
		if isTrain || (!isTrain && !isTest) {
			imagesBuf = bytes.NewReader(db.TrainImages)
			labelsBuf = bytes.NewReader(db.TrainLabels)
		} else if isTest {
			imagesBuf = bytes.NewReader(db.TestImages)
			labelsBuf = bytes.NewReader(db.TestLabels)
		}

		idx, err := mnistidx.NewIDX(imagesBuf, labelsBuf)
		if err != nil {
			log.Fatalln("mnistidx creation err:", err)
		}

		buf := idx.NewBuffer()

		i := 0

		for {
			label, err := idx.Read(buf)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			if imageIndex == -1 || i == imageIndex {
				f, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("label_%d_index_%d.png", label, i)))
				if err != nil {
					log.Fatal(err)
				}

				img := buf.GetNRGBAImg(idx.ImagesHeader)
				if err := png.Encode(f, img); err != nil {
					log.Fatal(err)
				}
				f.Close()
			}

			i++
		}
	} else {
		var imagesFilename, labelsFilename string
		var imagesDumpDb, labelsDumpDb []byte

		// if train set or both omitted
		if isTrain || (!isTrain && !isTest) {
			imagesFilename = "train-images.idx3-ubyte"
			labelsFilename = "train-labels.idx1-ubyte"
			imagesDumpDb = db.TrainImages
			labelsDumpDb = db.TrainLabels
		} else if isTest {
			imagesFilename = "t10k-images.idx3-ubyte"
			labelsFilename = "t10k-labels.idx1-ubyte"
			imagesDumpDb = db.TestImages
			labelsDumpDb = db.TestLabels
		}

		imagesFile, err := os.Create(filepath.Join(outputDir, imagesFilename))
		if err != nil {
			log.Fatal(err)
		}

		labelsFile, err := os.Create(filepath.Join(outputDir, labelsFilename))
		if err != nil {
			log.Fatal(err)
		}

		if _, err := imagesFile.Write(imagesDumpDb); err != nil {
			log.Fatal(err)
		}

		if _, err := labelsFile.Write(labelsDumpDb); err != nil {
			log.Fatal(err)
		}
	}
}
