package crumbjoiner

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"encoding/base64"
	"github.com/cheggaaa/pb/v3"
	"sort"
	proto "github.com/golang/protobuf/proto"
)

func Joiner(dirname string) {
	var fullfile strings.Builder
	filename := ""
	files, err := FilePathWalkDir(dirname)
	if err != nil {
		return
	}
	crumbs := []Crumb{}
	fmt.Println("Reading the Crumbs...")
	count := len(files)
	bar := pb.StartNew(count)
	for _, doc := range files {
		content, err := ioutil.ReadFile(doc)
		if err != nil {
			log.Fatal(err)
		}
		c := &Crumb{}
		err = proto.Unmarshal(content,c)
		if err != nil {
			log.Fatal(err)
		}
		crumbs = append(crumbs, *c)
		bar.Increment()
	}
	bar.Finish()
	fmt.Println("Joining the files")
	count = len(files)
	bar = pb.StartNew(count)
	filename = crumbs[0].Name
	sort.Slice(crumbs,func(i, j int) bool { return crumbs[i].Index < crumbs[j].Index })
	for _,crumb := range crumbs {
		fullfile.WriteString(crumb.Content)
		bar.Increment()
	}
	bar.Finish()
	// Converting from base64 to normal
	// Writing to the output file
	data, _ := base64.StdEncoding.DecodeString(fullfile.String())
	_ = ioutil.WriteFile(filename, []byte(data), 0644)
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
