package crumbjoiner

import (
 	"crypto/sha256"
	"fmt"
	"io/ioutil"
  "sync"
	"log"
	"github.com/cheggaaa/pb/v3"
	"sort"
  "os"
  "time"
  proto "github.com/golang/protobuf/proto"
)

func Integrity(dirname string) {
  var wg sync.WaitGroup
	files, err := FilePathWalkDir(dirname)
	if err != nil {
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
			err = proto.Unmarshal(content,c)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Lock()
			crumbs = append(crumbs, *c)
			mutex.Unlock()
			bar.Increment()
		}(doc)
		time.Sleep(10*time.Millisecond)
	}
	wg.Wait()
  bar.Finish()
	fmt.Println("Blockchain Verification Under Progress .....")
	sort.Slice(crumbs,func(i, j int) bool { return crumbs[i].Index < crumbs[j].Index })
	bar = pb.StartNew(len(crumbs))
	for i,crumb := range crumbs {
		if i==0{
			bar.Increment()

		}else{
      wg.Add(1)
      go func(prev_crum Crumb,crumb Crumb) {
        defer wg.Done()
        current_h := sha256.New()
  			c_prev_hash := sha256.New()
  			c_prev_hash.Write([]byte(prev_crum.Content))
  			c_prev_hash.Write([]byte(prev_crum.PrevHash))
  			current_h.Write([]byte(crumb.Content))
  			current_h.Write([]byte(c_prev_hash.Sum(nil)))
  			//fmt.Printf("%v %x %x\n",crumb.Index,current_h.Sum(nil),crumb.Hash)
  			if string(current_h.Sum(nil))==string(crumb.Hash) {
  				//fmt.Println("PIECE ",i)
  			}else{
  				fmt.Println("Blockchain Verification Failed!!")
  				os.Exit(1)
  			}
        bar.Increment()
      }(crumbs[i-1],crumb)
		}
	}
  wg.Wait()
	bar.Finish()
}
