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
	_ "strings"
	"io"
	"github.com/golang/protobuf/proto"
)

func Crumber(name string,size int) {
	//To chuck the file into 1 MB
	chuck := 1024 * 1024 * size
	//var content strings.Builder
	f, err := os.Open("./" + name)
	if err != nil {
		fmt.Println("File Not Found!")
		return
	}
	os.Mkdir(name+"-crumb", 0755)
	reader := bufio.NewReader(f)
	buf := make([]byte,chuck)
	var prev_crum Crumb
  for i := 0; ; i++ {
        n, err := reader.Read(buf)
        if err != nil {
           if err != io.EOF {
               log.Fatal(err)
           }
           break
        }
        //fmt.Print(string(buf[0:n]))
        //content.WriteString(string(buf[0:n]))
				loc := new(Crumb)
				if i==0 {
					loc.PrevHash = []byte("0000000")
				}else{
					loc.PrevHash = prev_crum.Hash
				}
				h := sha256.New()
				h.Write([]byte(buf[0:n]))
				h.Write([]byte(loc.PrevHash))
				loc.Index = int64(i)
				loc.Name = name
				loc.Hash = h.Sum(nil)
				loc.Content = base64.StdEncoding.EncodeToString([]byte(buf[0:n]))
				data := &Crumb{Index:loc.Index, Name:loc.Name, Hash:loc.Hash, PrevHash:loc.PrevHash, Content:loc.Content}
				b,err := proto.Marshal(data)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(name+"-crumb/"+name+strconv.Itoa(i)+".cbc", b, 0644)
				prev_crum = *data
  }
}
