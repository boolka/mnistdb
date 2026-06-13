# mnistdb

mnistdb provides embedded MNIST dataset assets and a small CLI helper to
extract images or dump the original IDX dataset files.

Features
- Embed the official MNIST IDX files into a Go package for easy access.
- CLI tool `extractor` to extract PNG images or dump the raw IDX files.

Background
The package embeds the original MNIST files (see https://yann.lecun.com/exdb/mnist/)
so you can include the datasets directly in Go binaries without shipping
separate files.

Contents
- `pkg/mnistdb`: provides four byte slices exported as `TrainImages`,
	`TrainLabels`, `TestImages`, `TestLabels` which are populated via `//go:embed`.
- `cmd/extractor`: a small command-line helper that can extract PNG images
	from the embedded datasets or write the raw IDX files to disk.

Requirements
- Go 1.16+ (for `embed` support)

Build
To build the CLI tool run:

```bash
go build -o extractor ./cmd/extractor
```

Usage
The `extractor` tool supports the following flags:

- `--train`        Use the train dataset (default)
- `--test`         Use the test dataset
- `--png`          Extract PNG digit images instead of dumping raw IDX files
- `--dir`, `-d`    Directory to output files (must exist)
- `--index`, `-i`  Extract only the image with the given index (zero-based)
- `--help`, `-h`   Print help

Examples

Extract all test images to `test` as PNGs:

```bash
./extractor --test --png -d test
```

Dump the train IDX files to the `train` directory:

```bash
./extractor --train -d train
```

Extract image with index 1000 from the train dataset to `train`:

```bash
./extractor --train --png -i 1000 -d train
```

Notes
- The datasets are embedded in [pkg/mnistdb](pkg/mnistdb/mnistdb.go).
- The `extractor` includes its usage text in `cmd/extractor/help.txt`.

Contributing
Feel free to open PRs to fix bugs, improve the CLI, or add more utilities
for working with IDX files.

License
This project is provided under the terms in the repository `LICENSE` file.
