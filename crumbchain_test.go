package main

import (
  "testing"
  "log"
  "os"
  "crypto/sha256"
  "os/exec"
  "io"
  "github.com/iamstefin/crumbchain/crumb"
  "github.com/iamstefin/crumbchain/crumbjoiner"
)

var hash = new(string);
var new_hash = new(string);

func TestCrumb(t *testing.T)  {
  cmd := exec.Command("mkfile", "-n","2g","test_file")
  err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
  crumb.Crumber("test_file",10)
  *hash = calchash("test_file")
}

func TestJoin(t *testing.T)  {
  cmd := exec.Command("rm","-rf","test_file")
  err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
  crumbjoiner.Joiner("test_file-crumb")
  cmd = exec.Command("rm","-rf","test_file-crumb")
  cmd.Run()
  *new_hash = calchash("test_file")

  if *hash == *new_hash {
    log.Printf("Done")
  }else{
    panic("Eerrr")
  }
  cmd = exec.Command("rm","-rf","test_file")
  err = cmd.Run()
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

  return(string(h.Sum(nil)))
}
