package main

import (
	"crypto/sha256"
	"github.com/iamstefin/crumbchain/crumb"
	"github.com/iamstefin/crumbchain/crumbjoiner"
	"io"
	"log"
	"os"
	"testing"
)

var hash = new(string)
var new_hash = new(string)

func TestCrumb(t *testing.T) {
	makefile("test_file", 500)
	crumb.Crumber("test_file", 10)
	*hash = calchash("test_file")
	err := os.Remove("test_file")
	if err != nil {
		panic("Error")
	}
}

func TestJoin(t *testing.T) {
	crumbjoiner.Joiner("test_file-crumb")
	*new_hash = calchash("test_file")

	if *hash == *new_hash {
		log.Printf("Done")
	} else {
		panic("Eerrr")
	}
}

func calchash(name string) string {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return (string(h.Sum(nil)))
}

func makefile(name string, s int) {
	size := int64(s * 1024 * 1024)
	fd, err := os.Create(name)
	if err != nil {
		log.Fatal("Failed to create output")
	}
	_, err = fd.Seek(size-1, 0)
	if err != nil {
		log.Fatal("Failed to seek")
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		log.Fatal("Write failed")
	}
	err = fd.Close()
	if err != nil {
		log.Fatal("Failed to close file")
	}
}
