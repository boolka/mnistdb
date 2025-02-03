// # mnistdb
//
// package provides direct url to download mnistdb database files (from [Github])
//
// Consts with database names:
//   - TrainImagesDb
//   - TrainLabelsDb
//   - TestImagesDb
//   - TestLabelsDb
//
// Join `MnistDbUrl` const with chosen database. For example
//
//	url.JoinPath(mnistdb.MnistDbUrl, mnistdb.TestImagesDb)
//
// to download test images database
//
// [Github]: https://github.com/
package mnistdb

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// MnistDbUrl the root URL to join to load the databases
const MnistDbUrl = "https://raw.githubusercontent.com/boolka/mnistdb/refs/heads/main/cmd/mnistcli/db/"

const (
	TrainImagesDb = "train-images.idx3-ubyte" // TrainImagesDb is train images database filename
	TrainLabelsDb = "train-labels.idx1-ubyte" // TrainLabelsDb is train labels database filename
	TestImagesDb  = "t10k-images.idx3-ubyte"  // TestImagesDb is test images database filename
	TestLabelsDb  = "t10k-labels.idx1-ubyte"  // TestLabelsDb is test labels database filename
)

// MnistDb structure to download mnist database files
type MnistDb struct {
	Dir string
}

// NewMnistDb create new [github.com/boolka/mnistdb/pkg/mnistdb.MnistDb] instance
func NewMnistDb(dir string) (*MnistDb, error) {
	if dir == "" {
		f, err := os.Stat(dir)

		if err != nil && os.IsNotExist(err) {
			return nil, err
		}

		if !f.IsDir() {
			return nil, errors.New(dir + " is not directory")
		}
	}

	return &MnistDb{
		Dir: dir,
	}, nil
}

// CheckDb to check is concrete database downloaded
func (m *MnistDb) CheckDb(dbName string) bool {
	_, err := os.Stat(filepath.Join(m.Dir, dbName))

	return err == nil
}

// UploadDb to upload database by name
func (m *MnistDb) UploadDb(dbName string) error {
	path, err := url.JoinPath(MnistDbUrl, dbName)

	if err != nil {
		return err
	}

	res, err := http.Get(path)

	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(m.Dir, dbName))

	if err != nil {
		return err
	}

	if _, err = io.Copy(f, res.Body); err != nil {
		return err
	}

	return nil
}

// GetDbPath to get database filepath by name
func (m *MnistDb) GetDbPath(dbName string) string {
	return filepath.Join(m.Dir, dbName)
}

// UploadMnistDbs to download all databases
func (m *MnistDb) UploadMnistDbs() error {
	for _, dbName := range []string{TrainImagesDb, TrainLabelsDb, TestImagesDb, TestLabelsDb} {
		isExists := m.CheckDb(dbName)

		if !isExists {
			if err := m.UploadDb(dbName); err != nil {
				return err
			}
		}
	}

	return nil
}
