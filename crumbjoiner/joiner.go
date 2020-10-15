package crumbjoiner

import (
	"encoding/base64"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	proto "github.com/golang/protobuf/proto"
)

func Joiner(dirname string) {
	var wg sync.WaitGroup
	filename := ""
	files, err := FilePathWalkDir(dirname)
	if err != nil {
		fmt.Println("Some error : ", err.Error())
		return
	}
	crumbs := []Crumb{}
	var mutex = &sync.Mutex{}
	fmt.Println("Reading the Crumbs...")
	count := len(files)
	bar := pb.StartNew(count)
	for _, doc := range files {
		wg.Add(1)
		go func(doc string) {
			defer wg.Done()
			content, err := ioutil.ReadFile(doc)
			if err != nil {
				log.Fatal(err)
			}
			c := &Crumb{}
			err = proto.Unmarshal(content, c)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Lock()
			crumbs = append(crumbs, *c)
			mutex.Unlock()
			bar.Increment()
		}(doc)
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
	bar.Finish()
	fmt.Println(len(crumbs))
	fmt.Println("Joining and Writing the files")
	count = len(files)
	bar = pb.StartNew(count)
	filename = crumbs[0].Name
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sort.Slice(crumbs, func(i, j int) bool { return crumbs[i].Index < crumbs[j].Index })
	for _, crumb := range crumbs {
		//fmt.Println(crumb.Content)
		tmp, err := base64.StdEncoding.DecodeString(string(crumb.Content))
		if err != nil {
			panic(err)
		}
		if _, err = f.WriteString(string(tmp)); err != nil {
			panic(err)
		}
		bar.Increment()
		crumb = Crumb{}
	}
	bar.Finish()
}

func FilePathWalkDir(root string) ([]string, error) {
	//Listing all the files in a directory
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
