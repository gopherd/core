package ioutil

import (
	"io"
	"os"
)

func ReadFromFile(r io.ReaderFrom, filename string) (n int64, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	r.ReadFrom(f)
	return
}

func WriteToFile(w io.WriterTo, filename string) (n int64, err error) {
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return w.WriteTo(f)
}
