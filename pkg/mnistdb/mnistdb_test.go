package mnistdb_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/boolka/mnistdb/pkg/mnistdb"
)

func TestGithubUrlsAvailability(t *testing.T) {
	t.Parallel()

	for _, dbName := range []string{
		mnistdb.TrainImagesDb,
		mnistdb.TrainLabelsDb,
		mnistdb.TestImagesDb,
		mnistdb.TestLabelsDb,
	} {
		path, err := url.JoinPath(mnistdb.MnistDbUrl, dbName)

		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Get(path)

		if err != nil || resp.StatusCode != 200 {
			t.Fatal(err)
		}
	}
}

func TestExistDbs(t *testing.T) {
	t.Parallel()

	mdb, err := mnistdb.NewMnistDb("../../cmd/mnistcli/db")

	if err != nil {
		t.Fatal(err)
	}

	for _, dbName := range []string{mnistdb.TrainImagesDb, mnistdb.TrainLabelsDb, mnistdb.TestImagesDb, mnistdb.TestLabelsDb} {
		isExists := mdb.CheckDb(dbName)

		if !isExists {
			t.Fatal(err)
		}
	}
}
