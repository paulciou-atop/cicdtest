package file

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

type File struct {
	name   string
	buffer *bytes.Buffer
}

func NewFile(name string) *File {
	return &File{
		name:   name,
		buffer: &bytes.Buffer{},
	}
}

func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)

	return err
}

func (f *File) Save() error {
	if err := ioutil.WriteFile(f.name, f.buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func Remove(path string) error {
	e := os.RemoveAll(path)
	if e != nil {
		log.Print(e)
	}
	return e
}
