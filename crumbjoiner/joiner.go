package crumbjoiner

import (
	"encoding/base64"
	"fmt"
	"github.com/francoispqt/gojay"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/cheggaaa/pb/v3"
	"sort"
)

type Crumb struct {
	index   int
	name    string
	hash    string
	prev_hash string
	content string
}

func (u *Crumb) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "index":
		return dec.Int(&u.index)
	case "name":
		return dec.String(&u.name)
	case "hash":
		return dec.String(&u.hash)
	case "prev_hash":
		return dec.String(&u.prev_hash)
	case "content":
		return dec.String(&u.content)
	}
	return nil
}

func (u *Crumb) NKeys() int {
	return 5
}

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
		err = gojay.UnmarshalJSONObject(content, c)
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
	filename = crumbs[0].name
	sort.Slice(crumbs,func(i, j int) bool { return crumbs[i].index < crumbs[j].index })
	for _,crumb := range crumbs {
		fullfile.WriteString(crumb.content)
		bar.Increment()
	}
	bar.Finish()
	fmt.Println("Decoding from base64.....")
	// Converting from base64 to normal
	data, err := base64.StdEncoding.DecodeString(fullfile.String())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	// Writing to the output file
	err = ioutil.WriteFile(filename, []byte(data), 0644)
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