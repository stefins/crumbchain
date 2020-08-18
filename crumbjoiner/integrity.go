package crumbjoiner

import (
 	"crypto/sha256"
	"fmt"
	"github.com/francoispqt/gojay"
	"io/ioutil"
	"log"
	"github.com/cheggaaa/pb/v3"
	"sort"
	"os"
)

func Integrity(dirname string) {
	files, err := FilePathWalkDir(dirname)
	if err != nil {
		return
	}
	crumbs := []Crumb{}
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
	}
	fmt.Println("Blockchain Verification Under Progress .....")
	sort.Slice(crumbs,func(i, j int) bool { return crumbs[i].index < crumbs[j].index })
	bar := pb.StartNew(len(crumbs))
	for i,crumb := range crumbs {
		if i==0{

		}else{
			current_h := sha256.New()
			c_prev_hash := sha256.New()
			c_prev_hash.Write([]byte(crumbs[i-1].content))
			c_prev_hash.Write([]byte(crumbs[i-1].prev_hash))
			current_h.Write([]byte(crumb.content))
			current_h.Write([]byte(c_prev_hash.Sum(nil)))
			//fmt.Printf("%v %x %x\n",crumb.index,current_h.Sum(nil),c_prev_hash.Sum(nil))
			if string(current_h.Sum(nil))==crumb.hash {
				//fmt.Println("PIECE ",i)
			}else{
				fmt.Println("Blockchain Verification Failed!!")
				os.Exit(1)
			}
		}
		bar.Increment()
	}
	bar.Finish()
}
