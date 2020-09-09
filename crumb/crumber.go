package crumb

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"encoding/base64"
	"io"

	"github.com/cheggaaa/pb/v3"
	"github.com/golang/protobuf/proto"
)

func Crumber(name string,size int) {
	// default size of the crumb
	chuck := 1024 * 1024 * size
	info, err := os.Stat(name)
	if err != nil {
		panic(err)
	}
	count := int(info.Size()/(1024*1024))/size
	bar := pb.StartNew(count)
	f, err := os.Open("./" + name)
	if err != nil {
		fmt.Println("File Not Found!")
		return
	}
	os.Mkdir(name+"-crumb", 0755)
	reader := bufio.NewReader(f)
	buf := make([]byte,chuck)
	var prev_crum Crumb
	fmt.Println("Writing Crumbs....")
  for i := 0; ; i++ {
        n, err := reader.Read(buf)
        if err != nil {
           if err != io.EOF {
               log.Fatal(err)
           }
           break
        }
				loc := new(Crumb)
				if i==0 {
					loc.PrevHash = []byte("0000000")
				}else{
					loc.PrevHash = prev_crum.Hash
				}
				h := sha256.New()
				loc.Content = base64.StdEncoding.EncodeToString([]byte(buf[0:n]))
				h.Write([]byte(loc.Content))
				h.Write([]byte(loc.PrevHash))
				loc.Index = int64(i)
				loc.Name = name
				loc.Hash = h.Sum(nil)
				data := &Crumb{Index:loc.Index, Name:loc.Name, Hash:loc.Hash, PrevHash:loc.PrevHash, Content:loc.Content}
				b,err := proto.Marshal(data)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(name+"-crumb/"+name+strconv.Itoa(i)+".cbc", b, 0644)
				prev_crum = *data
				bar.Increment()
  }
	bar.Finish()
}
