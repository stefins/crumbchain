package crumbjoiner

import (
	"crypto/sha256"
	"fmt"
	"github.com/francoispqt/gojay"
	"io/ioutil"
	"log"
	"os"
	"github.com/cheggaaa/pb/v3"
)

func Integrity(dirname string) {
	files, err := FilePathWalkDir(dirname)
	if err != nil {
		return
	}
	cr := []Crumb{}
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
		cr = append(cr, *c)
	}
	fmt.Println("Blockchain Verification Under Progress .....")
	count := len(files)
	bar := pb.StartNew(count)
	defer bar.Finish()
	for i := 0; i < len(cr); i++ {
		for j := 0; j < len(cr); j++ {
			if cr[j].index == i {
				if cr[j].index == len(cr)-1 {
					h := sha256.New()
					if (j-1<0){
						fmt.Println("Verification Failed!!!")
						os.Exit(1)
					}
					if cr[j].index == 0{
						curr_prev_hash := "0000000"
						h.Write([]byte(cr[j].content))
						h.Write([]byte(curr_prev_hash))
						//fmt.Printf("%x %x\n",h.Sum(nil),curr_prev_hash)
					}else{
						curr_prev_hash := sha256.New()
						curr_prev_hash.Write([]byte(cr[j-1].content))
						curr_prev_hash.Write([]byte(cr[j-1].prev_hash))
						h.Write([]byte(cr[j].content))
						h.Write([]byte(curr_prev_hash.Sum(nil)))
						//fmt.Printf("%x %x\n",h.Sum(nil),curr_prev_hash.Sum(nil))
					}
					if string(h.Sum(nil)) == cr[j].hash {
						bar.Increment()
						//fmt.Println("PIECE " + strconv.Itoa(cr[j].index) + " VERIFIED")
					} else {
						fmt.Println("Verification Failed!!!")
						os.Exit(1)
					}
				} else {
					h := sha256.New()
					if (j-1<0){
						fmt.Println("Verification Failed!!!")
						os.Exit(1)
					}
					if cr[j].index == 0{
						curr_prev_hash := "0000000"
						h.Write([]byte(cr[j].content))
						h.Write([]byte(curr_prev_hash))
						//fmt.Printf("%x %x\n",h.Sum(nil),curr_prev_hash)
					}else{
						curr_prev_hash := sha256.New()
						curr_prev_hash.Write([]byte(cr[j-1].content))
						curr_prev_hash.Write([]byte(cr[j-1].prev_hash))
						h.Write([]byte(cr[j].content))
						h.Write([]byte(curr_prev_hash.Sum(nil)))
						//fmt.Printf("%x %x\n",h.Sum(nil),curr_prev_hash.Sum(nil))
					}
					if string(h.Sum(nil)) == cr[j].hash {
						bar.Increment()
						//fmt.Println("PIECE " + strconv.Itoa(cr[j].index) + " VERIFIED")
					} else {
						fmt.Println("Verification Failed!!!")
						os.Exit(1)
					}
				}
				break
			}

		}
	}
}
