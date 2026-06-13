package mnistdb

import _ "embed"

//go:embed train-images.idx3-ubyte
var TrainImages []byte

//go:embed train-labels.idx1-ubyte
var TrainLabels []byte

//go:embed t10k-images.idx3-ubyte
var TestImages []byte

//go:embed t10k-labels.idx1-ubyte
var TestLabels []byte
